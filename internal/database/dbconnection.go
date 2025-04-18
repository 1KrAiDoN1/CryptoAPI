package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type Storage struct {
	DB *pgx.Conn
}

func ConnectDB() (*Storage, error) {
	cfg := GetDBconfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// err := godotenv.Load("/Users/pavelvasilev/Desktop/CryptoAPI/internal/database/secretHash.env")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DB_username, cfg.DB_password, cfg.DB_host, cfg.DB_port, cfg.DB_name)
	db, err := pgx.Connect(ctx, dbConnStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return &Storage{DB: db}, nil

}

func GetDBconfig() DB {
	err := godotenv.Load("/Users/pavelvasilev/Desktop/CryptoAPI/internal/database/DB_Config.env")
	if err != nil {
		log.Println("Ошибка при чтении конфигурации базы данных")
	}
	DB_config_path := DB{DB_username: os.Getenv("DB_USER"), DB_password: os.Getenv("DB_PASSWORD"), DB_host: os.Getenv("DB_HOST"), DB_port: os.Getenv("DB_PORT"), DB_name: os.Getenv("DB_NAME")}
	return DB_config_path
}
