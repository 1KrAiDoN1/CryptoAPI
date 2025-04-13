package handler

import (
	"context"
	// "errors"
	// "fmt"
	"helloapp/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
)

// RequireAuth is an authentication middleware that validates JWT tokens
// @Summary Authentication middleware
// @Description Verifies user authentication using JWT access token or refresh token flow.
// @Security ApiKeyAuth
// @Tags Authentication
// @Param Cookie header string true "Access token" default(access_token=your_token_here)
// @Param Cookie header string true "Refresh token" default(refresh_token=your_token_here)
// @Success 200 {object} object "Authenticated successfully"
// @Failure 302 "Redirect to login page when not authenticated"
// @Failure 401 {object} object "Unauthorized - Invalid or expired tokens"
// @SecurityDefinitions
// @ApiKeyAuth:
//
//	type: apiKey
//	in: cookie
//	name: access_token
//
// Middleware flow:
// 1. Checks for valid access_token cookie first
// 2. If access_token is invalid/expired, checks for refresh_token
// 3. With valid refresh_token:
//   - Generates new tokens
//   - Updates cookies and DB
//   - Continues request with new auth context
//
// 4. With no valid tokens: redirects to /login
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var userID int
		jwtCookie, jwtErr := r.Cookie("access_token")
		var err error
		if jwtErr == nil && jwtCookie != nil && jwtCookie.Value != "" {
			userID, err = service.ParseToken(jwtCookie.Value)
			if err != nil {
				log.Printf("JWT validation failed: %v", err)
			} else if userID > 0 {
				ctx = context.WithValue(ctx, "userID", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		refresh_token, _ := r.Cookie("refresh_token")
		if refresh_token != nil {
			userID, _ = Get_UserID_By_Refresh_Token(refresh_token.Value)
			if userID > 0 {
				email, password, err := Get_UserData_fromDB(userID)
				if err != nil {
					log.Printf("Ошибка при получении почты и пароля пользователя : %v", err)
				}

				err = Remove_The_Old_Refresh_Token(userID)
				if err != nil {
					log.Printf("Ошибка при удалении старого токена : %v", err)
				}

				// 3. Обновляем оба токена
				new_JWToken, err := service.GenerateJWToken(email, password)
				if err != nil {
					log.Printf("Ошибка при получении new_access_token: %v", err)
				}
				new_Refresh_Token, err := service.GenerateRefreshToken()
				if err != nil {
					log.Printf("Ошибка при получении new_Refresh_token : %v", err)
				}

				// 4. Устанавливаем новые куки и refresh токен в бд
				SetAuthCookies(w, new_JWToken, new_Refresh_Token)
				Save_New_Refresh_token(userID, new_Refresh_Token)
				ctx1 := context.WithValue(ctx, "userID", userID)
				next.ServeHTTP(w, r.WithContext(ctx1))

			} else {
				log.Println("Пользователь не авторизован: UserID == 0")
				log.Println(http.StatusUnauthorized)
			}

		} else {
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}

	}
}

func Remove_The_Old_Refresh_Token(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()
	connStr := "postgres://postgres:admin@localhost:5432/registration"
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return err
	}
	defer db.Close(ctx)

	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err = db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
func Save_New_Refresh_token(userID int, New_Refresh_Token string) {
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
	_, err = db.Exec(ctx, query, userID, New_Refresh_Token, RefreshTokenExpiresAt)
	if err != nil {
		log.Printf("Ошибка вставки данных: %v", err)
		return
	}
}
func SetAuthCookies(w http.ResponseWriter, new_JWToken, new_Refresh_Token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    new_JWToken,
		Expires:  time.Now().Add(service.TokenTTL),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    new_Refresh_Token,
		Expires:  time.Now().Add(service.RefreshTokenTTL),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

func Get_UserData_fromDB(userID int) (password, email string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()
	connStr := "postgres://postgres:admin@localhost:5432/registration"
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return "", "", err
	}
	defer db.Close(ctx)

	query := `SELECT email, password FROM users WHERE id = $1`
	err = db.QueryRow(ctx, query, userID).Scan(&email, &password)
	if err != nil {
		return "", "", err
	}
	return email, password, nil
}

func Get_UserID_By_Refresh_Token(refresh_token string) (userID int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()
	connStr := "postgres://postgres:admin@localhost:5432/registration"
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return 0, err
	}
	defer db.Close(ctx)
	query := `SELECT user_id FROM refresh_tokens WHERE token = $1`
	err = db.QueryRow(ctx, query, refresh_token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
