package main

import (
	"encoding/json"
	"fmt"

	//"helloapp/coincap"
	"html/template"
	"io"

	//"log"
	"net/http"
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
	// fmt.Println(string(body))

	// coincapClient, err := coincap.NewClient(time.Second * 10)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// assets, err := coincapClient.GetAssets()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, asset := range assets {
	// 	fmt.Println(asset.CorrectPrint())
	// }

	// bitcoin, err := coincapClient.GetAsset("bitcoin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(bitcoin.CorrectPrint())

	//
	//
	//
	clientUser := &http.Client{}
	nameCrypto := "xrp"
	responseClient, err := clientUser.Get(GetURL(nameCrypto))
	if err != nil {
		fmt.Println("Ошибка в получении данных", err)
	}
	defer responseClient.Body.Close()

	bodyResponse, err := io.ReadAll(responseClient.Body)
	if err != nil {
		fmt.Println("Ошибка в чтении данных", err)
	}
	fmt.Println(string(bodyResponse))

	var xrp CoinsStruct

	err = json.Unmarshal(bodyResponse, &xrp)
	if err != nil {
		fmt.Println("Ошибка при дессериализации", err)
	}
	CoinStack := fmt.Sprintf("ID: %s, Название: %s, Цена: %s", xrp.Coin.ID, xrp.Coin.Name, xrp.Coin.PriceUsd)
	fmt.Println(CoinStack)
	fmt.Println(xrp.Coin.ID, xrp.Coin.PriceUsd)

	HandleFunc(xrp)

}
func HandleFunc(xrp CoinsStruct) {
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		showInfo(w, r, xrp) // Передаем данные о криптовалюте в обработчик
	})
	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8080", nil)
}

func showInfo(w http.ResponseWriter, r *http.Request, xrp CoinsStruct) {
	Example := CoinStruct{
		Name:     xrp.Coin.Name,
		PriceUsd: xrp.Coin.PriceUsd,
	}
	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, Example)
	if err != nil {
		fmt.Println("Ошибка при отображении данных", err)
	}
}

func GetURL(nameCrypto string) string {
	return fmt.Sprintf("https://api.coincap.io/v2/assets/%s", nameCrypto)
}

// func (d CoinStruct) PostInfoCrypto(w http.ResponseWriter, r *http.Request) {

// }

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
