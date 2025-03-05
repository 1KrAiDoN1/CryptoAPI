package database

import (
	"context"
	"crypto/sha1"
	"fmt"
	"helloapp/internal/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

// postgres://<username>:<password>@<host>:<port>/<database>
// <username> — ваше имя пользователя.
// <password> — ваш пароль.
// <host> — адрес сервера (например, localhost для локального подключения или IP-адрес сервера).
// <port> — порт, на котором работает PostgreSQL (по умолчанию это 5432).
// <database> — имя базы данных, к которой вы хотите подключиться.

func SendUserRegistrationData(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email") // получаем данные из полей ввода пользователем
	Password := r.FormValue("password")
	timeOfRegistration := time.Now()
	if email == "" || Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	} else {
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

		query := `INSERT INTO users (email, password, registeredat) VALUES ($1, $2, $3)`
		_, err = db.Exec(ctx, query, email, Hash(Password), timeOfRegistration)
		if err != nil {
			log.Printf("Ошибка вставки данных: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		log.Printf("Пользователь с почтой %s зарегистрирован", email)
		log.Printf("ID пользователя с почтой %s в базе данных: %d", email, service.GetUserIdFromDB(email, Password))
		token, _ := service.GenerateJWToken(email, Password)
		log.Print("token пользователя:", token)
		id_from_token, _ := service.ParseToken(token)
		log.Print("Полученный ID из токена: ", id_from_token)

	}

}

func Hash(Password string) string {
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

type SignUpUser struct {
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}
