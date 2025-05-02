package main

import (
	_ "helloapp/docs"
	"helloapp/internal/handler"
)

// @title Crypto Market API
// @version 1.0
// @description API server for Crypto Market

// @contact.url https://t.me/KrAiDoN
// @contact.email pavelvasilev24843@gmail.com

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name access_token
// @description JWT token in cookie

// @securityDefinitions.apikey RefreshToken
// @in cookie
// @name refresh_token
// @description Refresh token in cookie
func main() {

	handler.HandleFunc()

	// TODO:
	// config init
	// logger init
	// storage init
	// main goroutine init (with graceful shutdown or os.Exit)

}
