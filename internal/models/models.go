package models

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
