package coincap

import (
	"fmt"
)

type AssetsResponse struct {
	Assets    []Asset `json: "data"`
	Timestamp int64   `json: "timestamp"`
}

func (d Asset) CorrectPrint() string {
	return fmt.Sprintf("[ID] %s | [Rank] %s | [Symbol] %s | [Name] %s | [Price_USD] %s", d.ID, d.Rank, d.Symbol, d.Name, d.PriceUsd)
}

type AssetResponse struct {
	Asset     Asset `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

type Asset struct {
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
