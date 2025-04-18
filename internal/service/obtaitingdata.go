package service

import (
	"context"
	"encoding/json"
	"fmt"
	"helloapp/internal/models"
	"io"

	"helloapp/internal/database"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func GetUserEmailFromDB(userID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg := database.GetDBconfig()
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Println("Error connecting to database", err)
	}
	defer db.DB.Close(ctx)

	var email string
	query := "SELECT email FROM users WHERE id = $1"
	err = db.DB.QueryRow(ctx, query, userID).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

func GetTimeOfRegistrationFromDB(userID int) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg := database.GetDBconfig()
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Println("Error connecting to database", err)
	}
	defer db.DB.Close(ctx)

	var timeOfRegistration string
	query := "SELECT registeredat FROM users WHERE id = $1"
	err = db.DB.QueryRow(ctx, query, userID).Scan(&timeOfRegistration)
	if err != nil {
		return time.Time{}, err
	}
	// Удаляем лишнюю часть строки (всё, что после временной зоны)
	timeString := strings.Split(timeOfRegistration, " MSK")[0]

	// Указываем формат строки
	layout := "2006-01-02 15:04:05.999999 -0700"

	// Преобразуем строку в time.Time
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func GetCryptoID(cryptoName string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg := database.GetDBconfig()
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Println("Error connecting to database", err)
	}
	defer db.DB.Close(ctx)
	var id int
	query := "SELECT id FROM cryptocurrencies WHERE name = $1"
	err = db.DB.QueryRow(ctx, query, cryptoName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetCryptoName(cryptoID int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg := database.GetDBconfig()
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Println("Error connecting to database", err)
	}
	defer db.DB.Close(ctx)

	var cryptoName string
	query := "SELECT name FROM cryptocurrencies WHERE id = $1"
	err = db.DB.QueryRow(ctx, query, cryptoID).Scan(&cryptoName)
	if err != nil {
		return "", err
	}
	return cryptoName, nil
}

func GetFavoriteCoins(user_id int) ([]models.CoinStruct, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg := database.GetDBconfig()
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Println("Error connecting to database", err)
	}
	defer db.DB.Close(ctx)

	query := `SELECT crypto_id FROM user_favorites WHERE user_id = $1`
	rows, err := db.DB.Query(ctx, query, user_id)
	if err != nil {
		return nil, fmt.Errorf("query favorites error: %v", err)
	}
	defer rows.Close()

	var cryptoIDs []int
	for rows.Next() {
		var crypto_id int
		err := rows.Scan(&crypto_id)
		if err != nil {
			return nil, fmt.Errorf("scan crypto_id error: %v", err)
		}
		cryptoIDs = append(cryptoIDs, crypto_id)
	}
	var Crypto []models.CoinStruct
	var wg sync.WaitGroup
	resultChan := make(chan models.CoinStruct)
	errChan := make(chan error)
	sem := make(chan struct{}, 20)                  // Ограничение на 10 одновременных запросов
	rateLimiter := time.Tick(50 * time.Millisecond) // 10 запросов в секунду

	// Обработчик результатов
	go func() {
		for value := range resultChan {
			Crypto = append(Crypto, value)
		}
	}()

	// Запуск горутин для каждого запроса
	for _, value := range cryptoIDs {
		<-rateLimiter // Ожидаем следующий тик для соблюдения rate limit
		wg.Add(1)
		go func(crypto_id int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			cryptoName, err := GetCryptoName(crypto_id)
			if err != nil {
				errChan <- fmt.Errorf("ошибка получения имени для ID %d: %v", crypto_id, err)
				return
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodGet,
				fmt.Sprintf("https://rest.coincap.io/v3/assets/%s?apiKey=4bb21996e06f4c5b19a7e9d09913ad52831d47014079a04aaf672118a6758677", cryptoName), nil)
			if err != nil {
				errChan <- err
				return
			}

			client := &http.Client{}
			response, err := client.Do(req)
			if err != nil {
				errChan <- err
				return
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				errChan <- fmt.Errorf("status code %d для %s", response.StatusCode, cryptoName)
				return
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				errChan <- err
				return
			}

			var cryptoInfo models.CoinsStruct
			if err := json.Unmarshal(body, &cryptoInfo); err != nil {
				errChan <- err
				return
			}

			resultChan <- cryptoInfo.Coin
		}(value)
	}
	// Ожидание завершения и сбор ошибок
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	// Обработка ошибок
	var errors []string
	go func() {
		for err := range errChan {
			errors = append(errors, err.Error())
		}
	}()

	// Блокировка до завершения всех горутин
	wg.Wait()

	if len(errors) > 0 {
		return Crypto, fmt.Errorf("некоторые запросы завершились с ошибками: %v", strings.Join(errors, "; "))
	}

	return Crypto, nil
}
