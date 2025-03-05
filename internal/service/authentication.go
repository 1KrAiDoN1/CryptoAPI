package service

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

const (
	TokenTTL = 24 * time.Hour
)

func GenerateJWToken(email string, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        strconv.Itoa(GetUserIdFromDB(email, password)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(TokenTTL).Unix(),
	})
	err := godotenv.Load("/Users/pavelvasilev/Desktop/HTTP & SQL with Go/internal/database/secretHash.env")
	if err != nil {
		log.Fatal(err)
	}
	secretSignInKey := os.Getenv("SECRET_SIGNINKEY") // получаем значение из файла конфигурации
	return token.SignedString([]byte(secretSignInKey))
}

func ParseToken(access_token string) (int, error) {
	token, err := jwt.ParseWithClaims(access_token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, errors.New("unexpected signing method")
		}
		secretSignInKey := os.Getenv("SECRET_SIGNINKEY")
		return []byte(secretSignInKey), nil // func ParseToken принимает токен, структуру для хранения данных о токене (в данном случае &jwt.StandardClaims{})
		// и функцию, которая возвращает секретный ключ для проверки подписи токена.
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return 0, errors.New("unexpected signing method")
	}
	UserID, _ := strconv.Atoi(claims.Id)
	return UserID, nil
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
	err = db.QueryRow(ctx, query, email, HashToken(password)).Scan(&UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		return 0
	}
	return UserId

}

func GetUserEmailFromDB(userID int) (string, error) {
	ctx := context.Background()
	connStr := "postgres://postgres:admin@localhost:5432/registration"
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return "", err
	}
	defer db.Close(ctx)

	var email string
	query := "SELECT email FROM users WHERE id = $1"
	err = db.QueryRow(ctx, query, userID).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

func HashToken(Password string) string {
	hash := sha1.New()
	err := godotenv.Load("/Users/pavelvasilev/Desktop/HTTP & SQL with Go/internal/database/secretHash.env")
	if err != nil {
		log.Fatal(err)
	}
	secretString := os.Getenv("SECRET_STRING") // получаем значение из файла конфигурации
	_, err1 := hash.Write([]byte(Password))
	if err1 != nil {
		log.Println("Ошибка при шифровании пароля", err)
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(secretString)))

}

type SignInUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
