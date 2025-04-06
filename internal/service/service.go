package service

import (
	"context"
	"encoding/json"
	"fmt"
	"helloapp/internal/models"
	"io"
	"net/http"
	"sync"
	"time"
)

// r.URL.Path
// r — это объект типа *http.Request, который представляет HTTP-запрос. Он содержит информацию о запросе, включая путь URL, заголовки, параметры и т.д.
// r.URL.Path — это поле, которое содержит путь URL, по которому был выполнен запрос. Например, если пользователь запрашивает URL http://localhost:8080/crypto/bitcoin, то r.URL.Path будет равно /crypto/bitcoin.
// 2. strings.TrimPrefix
// strings.TrimPrefix — это функция из пакета strings, которая удаляет указанный префикс (в данном случае "/crypto/") из строки, если он присутствует в начале.
// Синтаксис функции: strings.TrimPrefix(s, prefix), где s — строка, из которой нужно удалить префикс, а prefix — префикс, который нужно удалить.

func GetCryptoDataByID(id string) (models.CoinStruct, error) {
	var wg sync.WaitGroup
	var result models.CoinsStruct
	resultChan := make(chan models.CoinsStruct)
	errChan := make(chan error)
	defer close(resultChan)
	defer close(errChan)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://rest.coincap.io/v3/assets/%s?apiKey=4bb21996e06f4c5b19a7e9d09913ad52831d47014079a04aaf672118a6758677", id), nil)
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
			errChan <- err
			return
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			errChan <- err
			return
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()
	// go func() {
	// 	wg.Wait()
	// }()

	select {
	case result := <-resultChan:
		return result.Coin, nil
	case err := <-errChan:
		return models.CoinStruct{}, err

	case <-time.After(35 * time.Second):
		return models.CoinStruct{}, fmt.Errorf("request timeout")
	}

}

func GetCryptoData() ([]models.CoinStruct, error) {
	resultChan := make(chan []models.CoinStruct)
	errChan := make(chan error)
	defer close(resultChan)
	defer close(errChan)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		client := &http.Client{Timeout: 30 * time.Second}
		// req, err := client.Get("https://api.coincap.io/v2/assets/")
		req, err := client.Get("https://rest.coincap.io/v3/assets?apiKey=4bb21996e06f4c5b19a7e9d09913ad52831d47014079a04aaf672118a6758677")
		if err != nil {
			errChan <- err
			return
		}
		defer req.Body.Close()

		// Проверяем статус код ДО чтения тела
		if req.StatusCode != http.StatusOK {
			// Читаем тело ошибки для более информативного сообщения
			errorBody, _ := io.ReadAll(req.Body)
			errChan <- fmt.Errorf("неверный статус код %d: %s", req.StatusCode, string(errorBody))
			return
		}
		body, err := io.ReadAll(req.Body)
		if err != nil {
			errChan <- err
		}
		var spisok models.SliceCrypto
		err = json.Unmarshal(body, &spisok)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- spisok.Crypto

	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(35 * time.Second):
		return []models.CoinStruct{}, fmt.Errorf("request timeout")

	}
}
