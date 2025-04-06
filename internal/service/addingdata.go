package service

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

func AddFavoriteCryptoDB(userID int, cryptoName string) error {
	ctx := context.Background()
	connStr := "postgres://postgres:admin@localhost:5432/registration"
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Println("Ошибка при подключении к базе данных", err)
		return err
	}
	defer db.Close(ctx)

	cryptoID, err := GetCryptoID(cryptoName)
	if err != nil {
		log.Println("Ошибка при получении ID криптовалюты", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	query := `INSERT INTO user_favorites (user_id, crypto_id) VALUES ($1, $2)`
	_, err = db.Exec(ctx, query, userID, cryptoID)
	if err != nil {
		log.Printf("Ошибка вставки данных: %v", err)
		return err
	}
	return nil

}
