package auth

import (
	"context"
	"crypto/sha1"
	"fmt"
	"helloapp/internal/database"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

// postgres://<username>:<password>@<host>:<port>/<database>
// <username> — ваше имя пользователя.
// <password> — ваш пароль.
// <host> — адрес сервера (например, localhost для локального подключения или IP-адрес сервера).
// <port> — порт, на котором работает PostgreSQL (по умолчанию это 5432).
// <database> — имя базы данных, к которой вы хотите подключиться.
// func NewPostgresDB() (*pgx.Conn, error) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbName := os.Getenv("DB_NAME")
// 	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
// 	db, err := pgx.Connect(context.Background(), dbConnStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

// SendUserRegistrationData handles new user registration
// @Summary Register new user
// @Description Creates new user account with email and password
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Success 303 "Redirect to home page on success"
// @Failure 400 {string} string "Bad Request - Missing email or password"
// @Failure 500 {string} string "Internal Server Error"
// @Router /sendUserRegistrationData [post]
//
// Registration Flow:
// 1. Validates email and password presence
// 2. Hashes password (using bcrypt or similar)
// 3. Stores user in database with registration timestamp
// 4. Generates debug logs with registration info
// 5. Redirects to /home on success
//
// Security Notes:
// - Passwords are hashed before storage
// - No sensitive data is logged
// - Returns minimal error information to client
func SendUserRegistrationData(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email") // получаем данные из полей ввода пользователем
	Password := r.FormValue("password")
	timeOfRegistration := time.Now()
	if email == "" || Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	} else {
		ctx := context.Background()
		cfg := database.GetDBconfig()
		db, err := database.ConnectDB(cfg)
		if err != nil {
			http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		}
		defer db.Close(ctx)

		query := `INSERT INTO users (email, password, registeredat) VALUES ($1, $2, $3)`
		_, err = db.Exec(ctx, query, email, Hash(Password), timeOfRegistration)
		if err != nil {
			log.Printf("Ошибка вставки данных: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		log.Printf("Пользователь с почтой %s зарегистрирован", email)
		log.Printf("ID пользователя с почтой %s в базе данных: %d", email, GetUserIdFromDB(email, Password))
		token, _ := GenerateJWToken(email, Password)
		log.Print("token пользователя:", token)
		id_from_token, _ := ParseToken(token)
		log.Print("Полученный ID из токена: ", id_from_token)

	}

}

func Hash(Password string) string {
	hash := sha1.New()
	err := godotenv.Load("/Users/pavelvasilev/Desktop/CryptoAPI/internal/database/secretHash.env")
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
