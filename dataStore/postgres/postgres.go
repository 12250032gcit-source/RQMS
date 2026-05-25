package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	host := getEnv("DB_HOST", "dpg-d8a5h7ek1jcs73fl38q0-a.singapore-postgres.render.com")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "rqms_user")
	password := getEnv("DB_PASSWORD", "eFxzVEyP1owIQXM7gODMI7JUgoMLTOvt")
	dbname := getEnv("DB_NAME", "rqms")

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	Db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	} else {
		log.Println("Database connected successfully")
	}

	createTables()
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS signup (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(100),
			last_name  VARCHAR(100),
			email      VARCHAR(200) UNIQUE NOT NULL,
			password   VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS staff (
			id         SERIAL PRIMARY KEY,
			first_name VARCHAR(100),
			last_name  VARCHAR(100),
			email      VARCHAR(200) UNIQUE NOT NULL,
			password   VARCHAR(255) NOT NULL,
			role       VARCHAR(50) DEFAULT 'staff',
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id         SERIAL PRIMARY KEY,
			first_name VARCHAR(100),
			last_name  VARCHAR(100),
			email      VARCHAR(200),
			phone      VARCHAR(50),
			note       TEXT,
			time       VARCHAR(100),
			status     VARCHAR(50) DEFAULT 'waiting',
			table_no   VARCHAR(20) DEFAULT '',
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS tables (
			id        SERIAL PRIMARY KEY,
			table_no  VARCHAR(20) UNIQUE NOT NULL,
			capacity  INT DEFAULT 4,
			status    VARCHAR(50) DEFAULT 'available'
		)`,
	}

	for _, q := range queries {
		if _, err := Db.Exec(q); err != nil {
			log.Printf("Warning creating table: %v", err)
		}
	}

	// Add missing columns to existing users table (for upgrades)
	Db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'waiting'")
	Db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS table_no VARCHAR(20) DEFAULT ''")
	Db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW()")

	var count int
	Db.QueryRow("SELECT COUNT(*) FROM tables").Scan(&count)
	if count == 0 {
		for i := 1; i <= 10; i++ {
			Db.Exec("INSERT INTO tables (table_no, capacity) VALUES ($1, $2)",
				fmt.Sprintf("T%02d", i), 4)
		}
		log.Println("Seeded 10 restaurant tables")
	}
}
