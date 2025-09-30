// envloader.go
package config_manager

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads required environment variables and returns them in a struct
// You can extend this with as many env variables as you need
type DBEnvConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func InitEnv() {
	// Load .env file from project root (if present)
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, falling back to system env")
	}
}

// LoadDBEnv tries to read variables from environment and logs fatal if missing
func LoadDBEnv() DBEnvConfig {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		log.Fatal("missing environment variable: DB_HOST")
	}

	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		log.Fatal("missing environment variable: DB_PORT")
	}
	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		log.Fatal("missing environment variable: DB_USER")
	}

	pass, ok := os.LookupEnv("DB_PASS")
	if !ok {
		log.Fatal("missing environment variable: DB_PASS")
	}

	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatal("missing environment variable: DB_NAME")
	}

	return DBEnvConfig{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
		Name: dbname,
	}
}
