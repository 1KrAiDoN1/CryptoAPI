package handler

import (
	"context"
	"fmt"
	"helloapp/internal/database"
	"helloapp/internal/models"
	"helloapp/internal/service"
	"helloapp/pkg/format"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	//"database/sql"

	"github.com/jackc/pgx/v4"
)

func HandleFunc() {
	http.HandleFunc("/home", service.RequireAuth(showInfo)) // Передаем данные о криптовалюте в обработчик
	http.HandleFunc("/crypto/", service.RequireAuth(showCryptoDetails))
	http.HandleFunc("/sign_up", registration_window)
	http.HandleFunc("/login", authorization_window)
	http.HandleFunc("/verification", Verification_User)
	http.HandleFunc("/sendUserRegistrationData", database.SendUserRegistrationData)
	http.HandleFunc("/logout", logout)
	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8080", nil)

}

func logout(w http.ResponseWriter, r *http.Request) {
	// refreshCookie, err := r.Cookie("refresh_token")
	// if err != nil {
	// 	http.Redirect(w, r, "home", http.StatusSeeOther)
	// 	return
	// }

	// // Удаляем refresh токен из базы данных
	// ctx := context.Background()
	// connStr := "postgres://postgres:admin@localhost:5432/registration"
	// db, err := pgx.Connect(ctx, connStr)
	// if err != nil {
	// 	http.Error(w, "Database connection error", http.StatusInternalServerError)
	// 	return
	// }
	// defer db.Close(ctx)

	// _, err = db.Exec(ctx, "DELETE FROM refresh_tokens WHERE token = $1", refreshCookie.Value)
	// if err != nil {
	// 	http.Error(w, "Failed to delete refresh token", http.StatusInternalServerError)
	// 	return
	// }
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		//Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		//Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func registration_window(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../pkg/templates/registration.html")
	if err != nil {
		log.Print("Ошибка при чтении шаблона:", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Print("Ошибка при выполнении шаблона:", err)
	}

}

func authorization_window(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../pkg/templates/authorization.html")
	if err != nil {
		log.Print("Ошибка при чтении шаблона:", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Print("Ошибка при выполнении шаблона:", err)
	}
}

// func RefreshHandler(w http.ResponseWriter, r *http.Request) {
// 	refreshToken, err := r.Cookie("refresh_token")
// 	if err != nil {
// 		http.Error(w, "Refresh token required", http.StatusBadRequest)
// 		return
// 	}

// 	newJWToken, err := CheckRefreshTokenTTL(refreshToken.Value)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	// Обновляем JWT в куках
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "jwt_token",
// 		Value:    newJWToken,
// 		Expires:  time.Now().Add(service.TokenTTL),
// 		HttpOnly: true,
// 		Path:     "/",
// 	})

// 	w.WriteHeader(http.StatusOK)
// 	http.Redirect(w, r, "/home", http.StatusSeeOther)
// }

// ПРИ АВТОРИЗАЦИИ ВЫДАЕМ ПОЛЬЗОВАТЕЛЮ ДВА ТОКЕНА (JWT и REFRESH)
// УСТАНАВЛИВАЕМ ТОКЕНЫ В КУКИ
func Verification_User(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user service.SignInUser
		user.Email = r.FormValue("email")
		user.Password = r.FormValue("password")

		userID := service.GetUserIdFromDB(user.Email, user.Password)
		if userID == 0 {
			http.Error(w, "Invalid credentials", http.StatusBadRequest)
			return
		}

		//token, err := service.GenerateJWToken(user.Email, user.Password)
		// if err != nil {
		// 	http.Error(w, "Error generating token", http.StatusInternalServerError)
		// 	return
		// }
		JWToken, RefreshToken, err := service.GetTokens(user.Email, user.Password) // генерируем два токена
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		connStr := "postgres://postgres:admin@localhost:5432/registration"
		db, err := pgx.Connect(ctx, connStr)
		if err != nil {
			log.Fatalf("ошибка при коннекте к базе данных: %v\n", err)
		}
		defer db.Close(ctx)

		err = db.Ping(ctx)
		if err != nil {
			log.Fatal(err)
		}
		RefreshTokenExpiresAt := time.Now().Add(service.RefreshTokenTTL)
		query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
		_, err = db.Exec(ctx, query, userID, RefreshToken, RefreshTokenExpiresAt)
		if err != nil {
			log.Printf("Ошибка вставки данных: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Устанавливаем JWToken в куки
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token",
			Value:    JWToken,
			Expires:  time.Now().Add(service.TokenTTL),
			HttpOnly: true,
			//Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    RefreshToken,
			Expires:  time.Now().Add(service.RefreshTokenTTL),
			HttpOnly: true,
			//Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
}

type DataInfo struct {
	Email  string
	Output []models.CoinStruct
}

func showInfo(w http.ResponseWriter, r *http.Request) {
	output, err := service.GetCryptoData()
	if err != nil {
		fmt.Println("Ошибка при получении данных", err, http.StatusInternalServerError)
		return
	}
	// Получаем userID из контекста (если используется middleware)
	userIDVal := r.Context().Value("userID")
	var userID int
	if userIDVal != nil {
		var ok bool
		userID, ok = userIDVal.(int)
		if !ok {
			fmt.Println("Ошибка: userID в контексте имеет неверный тип")
			fmt.Println(userIDVal)
		}
	}
	var Email string
	if userID != 0 {
		// Получаем email пользователя из базы данных
		var err error
		Email, err = service.GetUserEmailFromDB(userID)
		if err != nil {
			fmt.Println("Ошибка при получении email пользователя", err)
			// Продолжаем выполнение без email
		}
	}
	data := DataInfo{Email: Email, Output: output}
	//fmt.Printf("Данные для шаблона: %+v\n", output)
	// Используем ParseFiles для загрузки шаблона из файла
	tmpl, err := template.New("home.html").Funcs(template.FuncMap{
		"formatLargeNumber": format.FormatLargeNumber, "formatLargeNumberForPercent": format.FormatLargeNumberForPercent, // Регистрируем функцию
	}).ParseFiles("../pkg/templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон с данными
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка при обработке шаблона", http.StatusInternalServerError)
		return
	}
}

func showCryptoDetails(w http.ResponseWriter, r *http.Request) {
	// получем userID из контекста
	// если токен не валиден или пользователь не авторизован, перенаправляем на страницу авторизации
	userIDVal := r.Context().Value("userID")
	// var userID int
	if userIDVal == nil || userIDVal == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//userID := r.Context().Value("userID").(int)
	// if userID == 0 {
	// 	http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	// 	return
	// }
	// Извлекаем ID криптовалюты из URL
	id := strings.TrimPrefix(r.URL.Path, "/crypto/")
	if id == "" {
		http.Error(w, "ID криптовалюты не указан", http.StatusBadRequest)
		return
	}

	// Получаем данные о конкретной криптовалюте
	crypto, err := service.GetCryptoDataByID(id)
	if err != nil {
		http.Error(w, "Ошибка при получении данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Загружаем шаблон для страницы с деталями
	tmpl, err := template.New("crypto_details.html").Funcs(template.FuncMap{
		"formatLargeNumber": format.FormatLargeNumber, "formatLargeNumberForPercent": format.FormatLargeNumberForPercent, // Регистрируем функцию
	}).ParseFiles("../pkg/templates/crypto_details.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон с данными
	err = tmpl.Execute(w, crypto)
	if err != nil {
		http.Error(w, "Ошибка при обработке шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
