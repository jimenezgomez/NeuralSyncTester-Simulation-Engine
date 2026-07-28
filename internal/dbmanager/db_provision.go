package dbmanager

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lib/pq"
)

// EnsureDatabase creates dbName if it does not already exist. It connects to
// Postgres's built-in "postgres" maintenance database to do so, since a
// database can't be created while connected to itself.
func EnsureDatabase(host, port, user, pass, dbName string) error {
	adminConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		user, pass, host, port)
	adminDB, err := sql.Open("postgres", adminConnString)
	if err != nil {
		return fmt.Errorf("open admin DB connection: %w", err)
	}
	defer adminDB.Close()

	var exists bool
	err = adminDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check database existence: %w", err)
	}

	if !exists {
		_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(dbName)))
		if err != nil {
			return fmt.Errorf("create database %q: %w", dbName, err)
		}
		log.Printf("Created database %q", dbName)
	}

	return nil
}

// ApplySchema executes every *.sql file in schemaDir against db. Schema files
// use CREATE TABLE IF NOT EXISTS, so this is safe to call on every startup.
func ApplySchema(db *sql.DB, schemaDir string) error {
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		return fmt.Errorf("read schema dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		path := filepath.Join(schemaDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read schema file %s: %w", path, err)
		}

		if _, err := db.Exec(string(data)); err != nil {
			return fmt.Errorf("apply schema file %s: %w", path, err)
		}
	}

	return nil
}
