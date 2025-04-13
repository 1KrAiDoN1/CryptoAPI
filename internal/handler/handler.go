package handler

import (
	"context"
	"fmt"
	"helloapp/internal/models"
	"helloapp/internal/service"
	"helloapp/pkg/format"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	//"database/sql"

	//"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func HandleFunc() {
	http.HandleFunc("/home", RequireAuth(showInfo)) // Передаем данные о криптовалюте в обработчик
	http.HandleFunc("/crypto/", RequireAuth(showCryptoDetails))
	http.HandleFunc("/personal_account", RequireAuth(showPersonalAccount))
	http.HandleFunc("/sign_up", registration_window)
	http.HandleFunc("/login", authorization_window)
	http.HandleFunc("/verification", Verification_User)
	http.HandleFunc("/saveFavoriteCrypto/", RequireAuth(saveFavoriteCrypto))
	http.HandleFunc("/sendUserRegistrationData", SendUserRegistrationData)
	http.HandleFunc("/logout", logout)
	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Явно указываем путь
	))

	log.Println("Сервер запущен")
	http.ListenAndServe(":8080", nil)

}

// saveFavoriteCrypto adds cryptocurrency to user's favorites
// @Summary Add favorite cryptocurrency
// @Description Saves specified cryptocurrency to authenticated user's favorites list
// @Tags User
// @Security ApiKeyAuth
// @Param crypto_name path string true "Cryptocurrency ID or symbol"
// @Success 303 "Redirect to personal account page"
// @Failure 400 {string} string "Bad Request - Missing cryptocurrency ID"
// @Failure 401 {string} string "Unauthorized - User not authenticated"
// @Failure 500 {string} string "Internal Server Error"
// @Router /saveFavoriteCrypto/{crypto_name} [post]
//
// Authentication:
// - Requires valid JWT token in cookies (via RequireAuth middleware)
// - Uses userID from request context
//
// Parameters:
// - crypto_name: Cryptocurrency identifier extracted from URL path
//
// Flow:
// 1. Verifies user authentication
// 2. Extracts cryptocurrency ID from URL path
// 3. Saves to user's favorites in database
// 4. Redirects to /personal_account
func saveFavoriteCrypto(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("userID").(int)
	if userIDVal == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Извлекаем ID криптовалюты из URL
	crypto_name := strings.TrimPrefix(r.URL.Path, "/saveFavoriteCrypto/")
	if crypto_name == "" {
		http.Error(w, "ID криптовалюты не указан", http.StatusBadRequest)
		return
	}
	err := service.AddFavoriteCryptoDB(userIDVal, crypto_name)
	if err != nil {
		log.Println("Ошибка при добавлении в избранное:", err)
	}
	http.Redirect(w, r, "/personal_account", http.StatusSeeOther)

}

// logout performs user logout procedure
// @Summary Logout user
// @Description Invalidates user session by clearing authentication cookies and removing refresh token
// @Tags Authentication
// @Produce json
// @Security ApiKeyAuth
// @Param Cookie header string true "Refresh token" default(refresh_token=your_token_here)
// @Success 303 "Redirect to login page"
// @Failure 401 {string} string "Unauthorized - No valid session"
// @Router /logout [post]
//
// Cookie Management:
// - Clears both access_token and refresh_token cookies
// - Sets expired cookies (past date) to ensure browser removal
// - Uses HttpOnly flag for security
//
// Database Operations:
// - Removes refresh token from database if exists
//
// Flow:
// 1. Checks for existing refresh token
// 2. Removes token from database if found
// 3. Clears both access and refresh token cookies
// 4. Redirects to login page
func logout(w http.ResponseWriter, r *http.Request) {
	refresh_token, _ := r.Cookie("refresh_token")
	if refresh_token != nil {
		userID, _ := Get_UserID_By_Refresh_Token(refresh_token.Value)
		Remove_The_Old_Refresh_Token(userID)
	} else {
		log.Println("userID == 0")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
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

type DataUser struct {
	Email              string
	TimeOfRegistration time.Time
	FavoriteCrypto     []models.CoinStruct
}

// showPersonalAccount displays user's personal account information
// @Summary Get personal account data
// @Description Retrieves authenticated user's profile information including email, registration time and favorite cryptocurrencies
// @Tags User
// @Produce html
// @Security ApiKeyAuth
// @Param Cookie header string true "Access token" default(access_token=your_token_here)
// @Success 200 {object} DataUser "HTML page with user data"
// @Failure 401 {string} string "Unauthorized - Missing or invalid token"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 504 {string} string "Gateway Timeout - Data loading timeout"
// @Router /personal_account [get]
//
// Response Structure (DataUser):
//
//	{
//	  "Email": "user@example.com",
//	  "TimeOfRegistration": "2023-01-15T10:30:00Z",
//	  "FavoriteCrypto": [
//	    {
//	      "ID": "bitcoin",
//	      "Symbol": "btc",
//	      "Name": "Bitcoin",
//	      ...
//	    }
//	  ]
//	}
//
// Notes:
// - Requires authentication via JWT cookie
// - Uses HTML templating with custom formatting functions
// - Implements 5-second timeout for loading favorite coins
// - May return partial data if some components fail to load
func showPersonalAccount(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(int)
	email, err := service.GetUserEmailFromDB(userID)
	if err != nil {
		log.Println("Ошибка при получении почты пользователя", err)
	}
	timeOfRegistration, err := service.GetTimeOfRegistrationFromDB(userID)
	if err != nil {
		log.Println("Ошибка при получении времени регистрации пользователя", err)
	}
	done := make(chan struct{})
	var favoriteCrypto []models.CoinStruct
	var loadErr error
	go func() {
		defer close(done)
		favoriteCrypto, loadErr = service.GetFavoriteCoins(userID)
	}()

	// Ожидание завершения с таймаутом
	select {
	case <-done:
		// Продолжаем выполнение
	case <-time.After(5 * time.Second):
		log.Println("Таймаут загрузки данных")
		http.Error(w, "Превышено время загрузки данных", http.StatusGatewayTimeout)
		return
	}

	if loadErr != nil {
		log.Printf("Ошибка загрузки данных: %v", loadErr)
		// Можно показать частичные данные
	}
	data := DataUser{Email: email, TimeOfRegistration: timeOfRegistration, FavoriteCrypto: favoriteCrypto}
	tmpl, err := template.New("personal_account.html").Funcs(template.FuncMap{"formatLargeNumber": format.FormatLargeNumber, "formatLargeNumberForPercent": format.FormatLargeNumberForPercent, "Float": format.Float}).ParseFiles("../pkg/templates/personal_account.html")
	if err != nil {
		log.Println("Ошибка при чтении шаблона", err)
		http.Error(w, "Ошибка при чтении шаблона", http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Print("Ошибка при выполнении шаблона:", err)
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
	}
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

// Verification_User authenticates user and issues JWT tokens
// @Summary User authentication
// @Description Verifies user credentials and issues access/refresh tokens
// @Tags Authentication
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Success 303 "Redirect to /home on success"
// @Success 200 {object} object "Sets access_token and refresh_token cookies"
// @Failure 401 {string} string "Unauthorized - Invalid credentials"
// @Failure 500 {string} string "Internal Server Error"
// @Router /verification [post]
//
// Cookies:
// - access_token: JWT access token (expires in TokenTTL)
// - refresh_token: Refresh token (expires in RefreshTokenTTL)
//
// Security:
// - Stores refresh token in database
// - Sets HttpOnly cookies for enhanced security
//
// Flow:
// 1. Validates email/password against database
// 2. Generates new JWT and refresh tokens
// 3. Stores refresh token in database
// 4. Sets secure HTTP cookies
// 5. Redirects to /home
func Verification_User(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user service.SignInUser
		user.Email = r.FormValue("email")
		user.Password = r.FormValue("password")

		userID := service.GetUserIdFromDB(user.Email, user.Password)
		if userID == 0 {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

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
			Name:     "access_token",
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
		//w.Header().Add("Access_token", JWToken)

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
}

type DataInfo struct {
	Email  string              `json: "email"`
	Output []models.CoinStruct `json: "data"`
}

// showInfo displays cryptocurrency data and user information
// @Summary Get cryptocurrency data and user info
// @Description Fetches cryptocurrency market data and displays it along with authenticated user's email (if available)
// @Tags Crypto
// @Produce html
// @Security ApiKeyAuth
// @Param Cookie header string true "Access token" default(access_token=your_token_here)
// @Success 200 {object} DataInfo "HTML page with crypto data"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 401 {string} string "Unauthorized (though page will still render without user data)"
// @Router /home [get]
//
// Response structure:
//
//	DataInfo {
//	   Email string `json:"email"`    // User's email (if authenticated)
//	   Output []service.CryptoData `json:"data"` // Array of cryptocurrency data
//	}
func showInfo(w http.ResponseWriter, r *http.Request) {
	output, err := service.GetCryptoData()
	if err != nil {
		fmt.Println("Ошибка при получении данных", err, http.StatusInternalServerError)
		return
	}

	//authHeader := r.Header.Get("Access_token")
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
	//fmt.Println(authHeader)

	// Выполняем шаблон с данными
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка при обработке шаблона", http.StatusInternalServerError)
		return
	}
}

// showCryptoDetails displays detailed information about a specific cryptocurrency
// @Summary Get cryptocurrency details
// @Description Shows detailed information page for specified cryptocurrency
// @Tags Crypto
// @Produce html
// @Security ApiKeyAuth
// @Param id path string true "Cryptocurrency ID"
// @Param Cookie header string true "Access token" default(access_token=your_token_here)
// @Success 200 "HTML page with cryptocurrency details"
// @Failure 302 "Redirect to login if not authenticated"
// @Failure 400 {string} string "Bad Request - Missing cryptocurrency ID"
// @Failure 404 {string} string "Not Found - Cryptocurrency not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /crypto/{id} [get]
//
// Template Functions:
// - formatLargeNumber: Formats large numbers for display
// - formatLargeNumberForPercent: Formats percentage values
//
// Flow:
// 1. Checks user authentication via context
// 2. Extracts cryptocurrency ID from URL path
// 3. Fetches detailed crypto data from service
// 4. Renders crypto_details.html template with data
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
