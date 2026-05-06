package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresDatabase(env *Env) *sql.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass
	dbName := env.DBName
	dbSSLMode := env.DBSSLMode

	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode)
	if dbUser == "" || dbPass == "" {
		postgresDSN = fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s", dbHost, dbPort, dbName, dbSSLMode)
	}

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := ensurePostgresSchema(ctx, db); err != nil {
		log.Fatal(err)
	}

	return db
}

func ClosePostgresDBConnection(db *sql.DB) {
	if db == nil {
		return
	}

	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to PostgreSQL closed.")
}

func ensurePostgresSchema(ctx context.Context, db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user'`,
		`CREATE TABLE IF NOT EXISTS trips (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			start_time BIGINT,
			end_time BIGINT,
			total_distance DOUBLE PRECISION DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`ALTER TABLE trips ADD COLUMN IF NOT EXISTS total_distance DOUBLE PRECISION DEFAULT 0`,
		`CREATE TABLE IF NOT EXISTS places (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			type TEXT,
			category TEXT,
			rating DOUBLE PRECISION,
			open_time BIGINT,
			closed_time BIGINT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS trip_routes (
			id TEXT PRIMARY KEY,
			trip_id TEXT NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
			place_id TEXT NOT NULL REFERENCES places(id) ON DELETE CASCADE,
			order_index BIGINT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS favorites (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			place_id TEXT NOT NULL REFERENCES places(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, place_id)
		)`,
	}

	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}

	return nil
}
