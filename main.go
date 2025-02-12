package main

import (
	"encoding/json"
	"fmt"

	"helloapp/coincap"
	"html/template"
	"io"

	"log"
	"net/http"
	"time"
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

	coincapClient, err := coincap.NewClient(time.Second * 10)
	if err != nil {
		log.Fatal(err)
	}

	bitcoin, err := coincapClient.GetAsset("bitcoin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bitcoin.CorrectPrint())

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
	tmpl, err := template.ParseFiles("home.html")
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

func getCryptoData() ([]CoinStruct, error) {
	var CryptoBase []CoinStruct
	for _, id := range SliceOfNameCrypto {
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
