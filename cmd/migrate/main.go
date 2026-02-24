package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	action := flag.String("action", "up", "Migration action (up, down)")
	flag.Parse()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/finance_dashboard?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	switch *action {
	case "up":
		if err := runMigrations(db, "up"); err != nil {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := runMigrations(db, "down"); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

func runMigrations(db *sql.DB, direction string) error {
	migrations := []struct {
		version int
		up      string
		down    string
	}{
		{
			version: 1,
			up: `
				CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
				
				CREATE TABLE users (
					id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
					email VARCHAR(255) NOT NULL UNIQUE,
					password_hash VARCHAR(255),
					google_id VARCHAR(255) UNIQUE,
					global_currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
					created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
				);
				
				CREATE INDEX idx_users_email ON users(email);
				CREATE INDEX idx_users_google_id ON users(google_id);
				
				CREATE OR REPLACE FUNCTION update_updated_at_column()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.updated_at = CURRENT_TIMESTAMP;
					RETURN NEW;
				END;
				$$ LANGUAGE plpgsql;
				
				CREATE TRIGGER update_users_updated_at
					BEFORE UPDATE ON users
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();
			`,
			down: "DROP TABLE IF EXISTS users; DROP FUNCTION IF EXISTS update_updated_at_column(); DROP EXTENSION IF EXISTS uuid-ossp;",
		},
		{
			version: 2,
			up: `
				CREATE TABLE categories (
					id SERIAL PRIMARY KEY,
					name VARCHAR(100) NOT NULL UNIQUE,
					is_default BOOLEAN NOT NULL DEFAULT false,
					created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
				);
				
				INSERT INTO categories (name, is_default) VALUES
					('Продукты', true),
					('Транспорт', true),
					('Рестораны', true),
					('Здоровье', true),
					('Развлечения', true),
					('Дом', true),
					('Одежда', true),
					('Красота', true),
					('Образование', true),
					('Переводы', true),
					('Налоги и сборы', true),
					('Доходы', true),
					('Другое', true);
			`,
			down: "DROP TABLE IF EXISTS categories;",
		},
		{
			version: 3,
			up: `
				CREATE TABLE transactions (
					id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
					user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
					amount DECIMAL(15, 2) NOT NULL,
					currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
					description TEXT NOT NULL,
					date TIMESTAMP WITH TIME ZONE NOT NULL,
					place_name VARCHAR(255),
					place_lat DECIMAL(10, 8),
					place_lon DECIMAL(11, 8),
					category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
					is_confirmed BOOLEAN NOT NULL DEFAULT false,
					created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
				);
				
				CREATE INDEX idx_transactions_user_id ON transactions(user_id);
				CREATE INDEX idx_transactions_date ON transactions(date);
				CREATE INDEX idx_transactions_category_id ON transactions(category_id);
				CREATE INDEX idx_transactions_user_date ON transactions(user_id, date);
				
				CREATE TRIGGER update_transactions_updated_at
					BEFORE UPDATE ON transactions
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();
			`,
			down: "DROP TABLE IF EXISTS transactions;",
		},
		{
			version: 4,
			up: `
				CREATE TABLE user_category_rules (
					id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
					user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
					keyword VARCHAR(255) NOT NULL,
					category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
					created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(user_id, keyword)
				);
				
				CREATE INDEX idx_user_category_rules_user_id ON user_category_rules(user_id);
				CREATE INDEX idx_user_category_rules_keyword ON user_category_rules(keyword);
			`,
			down: "DROP TABLE IF EXISTS user_category_rules;",
		},
	}

	if direction == "up" {
		for _, m := range migrations {
			log.Printf("Applying migration version %d...", m.version)
			if _, err := db.Exec(m.up); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", m.version, err)
			}
		}
	} else {
		for i := len(migrations) - 1; i >= 0; i-- {
			m := migrations[i]
			log.Printf("Rolling back migration version %d...", m.version)
			if _, err := db.Exec(m.down); err != nil {
				return fmt.Errorf("failed to rollback migration %d: %w", m.version, err)
			}
		}
	}

	return nil
}
