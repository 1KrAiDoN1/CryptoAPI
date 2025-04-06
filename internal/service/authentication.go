package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

const (
	TokenTTL        = 24 * time.Hour
	RefreshTokenTTL = 30 * 24 * time.Hour
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

func GetTokens(email string, password string) (string, string, error) {
	JWToken, err := GenerateJWToken(email, password)
	if err != nil {
		return "", "", err
	}
	RefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	return JWToken, RefreshToken, nil
}

func GenerateRefreshToken() (string, error) {
	refresh_token := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	if _, err := r.Read(refresh_token); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", refresh_token), nil

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
		return 0, fmt.Errorf("token parse error: %w", err)
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		userID, err := strconv.Atoi(claims.Id)
		if err != nil {
			return 0, fmt.Errorf("invalid user ID in token")
		}
		return userID, nil
	}

	return 0, errors.New("invalid token")
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
