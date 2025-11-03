package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env not found")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL not found")
	}

	config, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Fatal("Error on parse database configuration")
	}

	config.MaxConns = 10
	config.MinConns = 20
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		DB, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			err = DB.Ping(context.Background())
			if err == nil {
				break
			}
		}

		log.Printf("Attempt %d/%d: Waiting for database... (%v)", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Error connecting to database after multiple attempts:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL database")

	if os.Getenv("RUN_MIGRATION") == "true" {
		runMigrations()
	}
}

func runMigrations() {
	fmt.Println("📝 Running migrations...")

	migrationsDir := "./migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Printf("⚠️  Could not read migrations directory: %v", err)
		return
	}

	// Ordena os arquivos
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Executa cada migration
	for _, filename := range sqlFiles {
		filePath := filepath.Join(migrationsDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("❌ Error reading migration %s: %v", filename, err)
			continue
		}

		_, err = DB.Exec(context.Background(), string(content))
		if err != nil {
			log.Printf("❌ Error executing migration %s: %v", filename, err)
			continue
		}
		fmt.Printf("✅ Migration executed: %s\n", filename)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("🔒 Conexão com o banco de dados fechada")
	}
}
