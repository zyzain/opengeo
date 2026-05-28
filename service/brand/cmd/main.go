package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"opengeo/service/brand/internal/adapter/mysql"
	"opengeo/service/brand/internal/application"
)

func main() {
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "root")
	dbName := getEnv("DB_NAME", "opengeo")
	servicePort := getEnv("SERVICE_PORT", "8002")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	brandRepo := mysql.NewBrandRepository(db)
	metadataRepo := mysql.NewBrandMetadataRepository(db)
	glossaryRepo := mysql.NewGlossaryRepository(db)
	snapshotRepo := mysql.NewSnapshotRepository(db)

	brandService := application.NewBrandService(brandRepo, metadataRepo, glossaryRepo, snapshotRepo)
	_ = brandService

	listener, err := net.Listen("tcp", ":"+servicePort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	log.Printf("Brand service listening on :%s", servicePort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down brand service...")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
