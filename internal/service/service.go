package service

import (
	"context"
	"encoding/json"
	"fmt"
	"helloapp/internal/models"
	"io"
	"log"
	"net/http"
	"time"
)

// r.URL.Path
// r — это объект типа *http.Request, который представляет HTTP-запрос. Он содержит информацию о запросе, включая путь URL, заголовки, параметры и т.д.
// r.URL.Path — это поле, которое содержит путь URL, по которому был выполнен запрос. Например, если пользователь запрашивает URL http://localhost:8080/crypto/bitcoin, то r.URL.Path будет равно /crypto/bitcoin.
// 2. strings.TrimPrefix
// strings.TrimPrefix — это функция из пакета strings, которая удаляет указанный префикс (в данном случае "/crypto/") из строки, если он присутствует в начале.
// Синтаксис функции: strings.TrimPrefix(s, prefix), где s — строка, из которой нужно удалить префикс, а prefix — префикс, который нужно удалить.

func GetCryptoDataByID(id string) (models.CoinStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.coincap.io/v2/assets/%s", id), nil)
	if err != nil {
		return models.CoinStruct{}, fmt.Errorf("oшибка при запросе к api в контексте: %s", err)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	//response, err := client.Get(fmt.Sprintf("https://api.coincap.io/v2/assets/%s", id))
	if err != nil {
		return models.CoinStruct{}, fmt.Errorf("ошибка при запросе к API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return models.CoinStruct{}, fmt.Errorf("API вернул статус: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return models.CoinStruct{}, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	var result models.CoinsStruct
	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.CoinStruct{}, fmt.Errorf("ошибка при разборе JSON: %v", err)
	}

	return result.Coin, nil
}

// func GetCryptoData() ([]models.CoinStruct, error) {
// 	var CryptoBase []models.CoinStruct
// 	for _, id := range SliceOfNameCrypto[:5] {
// 		clientUser := &http.Client{}
// 		response, err := clientUser.Get(fmt.Sprintf("https://api.coincap.io/v2/assets/%s", id))
// 		if err != nil {
// 			return []models.CoinStruct{}, err
// 		}
// 		defer response.Body.Close()

// 		bodyResponse, err := io.ReadAll(response.Body)
// 		if err != nil {
// 			return []models.CoinStruct{}, err
// 		}

// 		var result models.CoinsStruct
// 		err = json.Unmarshal(bodyResponse, &result)
// 		if err != nil {
// 			return []models.CoinStruct{}, err
// 		}
// 		CryptoBase = append(CryptoBase, result.Coin)

// 	}
// 	return CryptoBase, nil

// }

// var SliceOfNameCrypto = GetSliceOfNameCrypto()

func GetCryptoData() ([]models.CoinStruct, error) {
	client := &http.Client{}
	r, err := client.Get("https://api.coincap.io/v2/assets/")
	if err != nil {
		log.Println("Ошибка при запросе данных о криптовалютах:", err)
		return []models.CoinStruct{}, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		log.Println("Ошибка при запросе данных о криптовалютах:", err)
		return []models.CoinStruct{}, err

	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Ошибка при чтении данных о криптовалютах:", err)
		return []models.CoinStruct{}, err
	}

	var spisok models.SliceCrypto
	err = json.Unmarshal(body, &spisok)
	if err != nil {
		log.Println("Ошибка при десереализации данных о криптовалютах:", err)
		return []models.CoinStruct{}, err
	}

	return spisok.Crypto, nil
}
