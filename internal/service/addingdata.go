package service

import (
	"context"
	"log"
	"time"

	"helloapp/internal/database"
)

func AddFavoriteCryptoDB(userID int, cryptoName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := database.ConnectDB()
	if err != nil {
		log.Println("Error connecting to database", err)
	}
	defer db.DB.Close(ctx)

	cryptoID, err := GetCryptoID(cryptoName)
	if err != nil {
		log.Println("Ошибка при получении ID криптовалюты", err)
	}

	query := `INSERT INTO user_favorites (user_id, crypto_id) VALUES ($1, $2)`
	_, err = db.DB.Exec(ctx, query, userID, cryptoID)
	if err != nil {
		log.Printf("Ошибка вставки данных: %v", err)
		return err
	}
	return nil

}
