package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

const (
	TokenLive = 24 * time.Hour
)

func GenerateJWToken(email string, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(GetUserIdFromDB(email, password)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(TokenLive).Unix(),
	})
	err := godotenv.Load("/Users/pavelvasilev/Desktop/HTTP & SQL with Go/internal/database/secretHash.env")
	if err != nil {
		log.Fatal(err)
	}
	secretSignInKey := os.Getenv("SECRET_SIGNINKEY") // получаем значение из файла конфигурации
	return token.SignedString([]byte(secretSignInKey))

}

func GetUserIdFromDB(email string, password string) int {
	ctx := context.Background()
	connStr := "postgres://postgres:admin@localhost:5432/registration"
	// подключение к базе данных
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	// Проверка подключения
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var UserId int
	query := "SELECT id FROM users WHERE email = $1 AND password = $2"
	err = db.QueryRow(ctx, query, email, Hash(password)).Scan(&UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		return 0
	}
	return UserId

}

type SignInUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
