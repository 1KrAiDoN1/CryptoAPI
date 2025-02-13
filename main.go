package main

import (
	"encoding/json"
	"fmt"

	//"helloapp/coincap"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"time"
)

func main() {

	// client := &http.Client{}
	// r, err := client.Get("https://api.coincap.io/v2/assets/")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer r.Body.Close()
	// if r.StatusCode != http.StatusOK {
	// 	log.Fatal(r.Status)
	// }
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var spisok SliceCrypto
	// // spisok := make([]string, 0)
	// // for _, v := range body {
	// // 	spisok = append(spisok, v)
	// // }
	// // fmt.Println(spisok)
	// err = json.Unmarshal(body, &spisok)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// //fmt.Println(spisok)
	// var partOfMarket []string
	// for _, value := range spisok.Crypto {
	// 	partOfMarket = append(partOfMarket, value.ID)
	// }
	// fmt.Println(partOfMarket)

	// fmt.Println(string(body))

	// coincapClient, err := coincap.NewClient(time.Second * 10)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// bitcoin, err := coincapClient.GetAsset("bitcoin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//fmt.Println(bitcoin.CorrectPrint())

	//
	//
	//

	//CoinStack := fmt.Sprintf("ID: %s, Название: %s, Цена: %s", xrp.Coin.ID, xrp.Coin.Name, xrp.Coin.PriceUsd)
	// fmt.Println(CoinStack)
	// fmt.Println(xrp.Coin.ID, xrp.Coin.PriceUsd)
	HandleFunc()

}
func HandleFunc() {
	http.HandleFunc("/home", showInfo) // Передаем данные о криптовалюте в обработчик
	http.HandleFunc("/crypto/", showCryptoDetails)
	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8080", nil)
}

func showInfo(w http.ResponseWriter, r *http.Request) {
	output, err := getCryptoData()
	if err != nil {
		fmt.Println("Ошибка при получении данных", err, http.StatusInternalServerError)
		return
	}

	//fmt.Printf("Данные для шаблона: %+v\n", output)
	// Используем ParseFiles для загрузки шаблона из файла
	tmpl, err := template.New("home.html").Funcs(template.FuncMap{
		"formatLargeNumber": formatLargeNumber, "formatLargeNumberForPercent": formatLargeNumberForPercent, // Регистрируем функцию
	}).ParseFiles("home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон с данными
	err = tmpl.Execute(w, output)
	if err != nil {
		http.Error(w, "Ошибка при обработке шаблона", http.StatusInternalServerError)
		return
	}
}

func showCryptoDetails(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID криптовалюты из URL
	id := strings.TrimPrefix(r.URL.Path, "/crypto/")
	if id == "" {
		http.Error(w, "ID криптовалюты не указан", http.StatusBadRequest)
		return
	}

	// Получаем данные о конкретной криптовалюте
	crypto, err := getCryptoDataByID(id)
	if err != nil {
		http.Error(w, "Ошибка при получении данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Загружаем шаблон для страницы с деталями
	tmpl, err := template.New("crypto_details.html").Funcs(template.FuncMap{
		"formatLargeNumber": formatLargeNumber, "formatLargeNumberForPercent": formatLargeNumberForPercent, // Регистрируем функцию
	}).ParseFiles("crypto_details.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон с данными
	err = tmpl.Execute(w, crypto)
	if err != nil {
		http.Error(w, "Ошибка при обработке шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// r.URL.Path
// r — это объект типа *http.Request, который представляет HTTP-запрос. Он содержит информацию о запросе, включая путь URL, заголовки, параметры и т.д.
// r.URL.Path — это поле, которое содержит путь URL, по которому был выполнен запрос. Например, если пользователь запрашивает URL http://localhost:8080/crypto/bitcoin, то r.URL.Path будет равно /crypto/bitcoin.
// 2. strings.TrimPrefix
// strings.TrimPrefix — это функция из пакета strings, которая удаляет указанный префикс (в данном случае "/crypto/") из строки, если он присутствует в начале.
// Синтаксис функции: strings.TrimPrefix(s, prefix), где s — строка, из которой нужно удалить префикс, а prefix — префикс, который нужно удалить.

func getCryptoDataByID(id string) (CoinStruct, error) {
	client := &http.Client{}
	response, err := client.Get(fmt.Sprintf("https://api.coincap.io/v2/assets/%s", id))
	if err != nil {
		return CoinStruct{}, fmt.Errorf("ошибка при запросе к API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return CoinStruct{}, fmt.Errorf("API вернул статус: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return CoinStruct{}, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	var result CoinsStruct
	err = json.Unmarshal(body, &result)
	if err != nil {
		return CoinStruct{}, fmt.Errorf("ошибка при разборе JSON: %v", err)
	}

	return result.Coin, nil
}

func getCryptoData() ([]CoinStruct, error) {
	var CryptoBase []CoinStruct
	for _, id := range SliceOfNameCrypto[:10] {
		clientUser := &http.Client{}
		response, err := clientUser.Get(fmt.Sprintf("https://api.coincap.io/v2/assets/%s", id))
		if err != nil {
			return []CoinStruct{}, err
		}
		defer response.Body.Close()

		bodyResponse, err := io.ReadAll(response.Body)
		if err != nil {
			return []CoinStruct{}, err
		}

		var result CoinsStruct
		err = json.Unmarshal(bodyResponse, &result)
		if err != nil {
			return []CoinStruct{}, err
		}
		CryptoBase = append(CryptoBase, result.Coin)

	}
	return CryptoBase, nil
}

var SliceOfNameCrypto = getSliceOfNameCrypto()

func getSliceOfNameCrypto() []string {
	client := &http.Client{}
	r, err := client.Get("https://api.coincap.io/v2/assets/")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		log.Fatal(r.Status)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var spisok SliceCrypto
	err = json.Unmarshal(body, &spisok)
	if err != nil {
		log.Fatal(err)
	}

	var partOfMarket []string
	for _, value := range spisok.Crypto {
		partOfMarket = append(partOfMarket, value.ID)
	}
	return partOfMarket
}

type SliceCrypto struct {
	Crypto []CoinStruct `json:"data"`
	Timet1 int64        `json:"timestamp"`
}
type CoinsStruct struct {
	Coin  CoinStruct `json:"data"`
	Timet int64      `json:"timestamp"`
}
type CoinStruct struct {
	ID                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
}

func formatLargeNumber(num string) string {
	// Парсим строку в число с плавающей точкой
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return "N/A" // Возвращаем "N/A", если число невалидно
	}

	// Форматируем число в зависимости от его величины
	switch {
	case f >= 1e12: // Триллионы
		return fmt.Sprintf("%.2fT", f/1e12)
	case f >= 1e9: // Миллиарды
		return fmt.Sprintf("%.2fB", f/1e9)
	case f >= 1e6: // Миллионы
		return fmt.Sprintf("%.2fM", f/1e6)
	default: // Меньше миллиона
		return fmt.Sprintf("%.2f", f)
	}
}

func formatLargeNumberForPercent(num string) string {
	// Парсим строку в число с плавающей точкой
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return "N/A" // Возвращаем "N/A", если число невалидно
	}

	// Форматируем число в зависимости от его величины
	switch {
	case f > 0:
		return fmt.Sprintf("+%.2f", f)
	default:
		return fmt.Sprintf("%.2f", f)
	}

}
