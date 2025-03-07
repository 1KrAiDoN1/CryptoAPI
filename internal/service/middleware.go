package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
)

// ПРОВЕРЯЕМ ТОКЕНЫ НА ВАЛИДНОСТЬ И
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из куки
		// Req_ctx := r.Context()
		// cookie, err := r.Cookie("jwt_token")
		// if err != nil {
		// 	//http.Redirect(w, r, "/sign_in", http.StatusUnauthorized)
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// // Парсим токен
		// userID, err := ParseToken(cookie.Value)
		// if err != nil || userID == 0 {
		// 	//http.Redirect(w, r, "/sign_in", http.StatusUnauthorized)
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// // Добавляем userID в контекст запроса
		// ctx := context.WithValue(Req_ctx, "userID", userID)
		// next.ServeHTTP(w, r.WithContext(ctx))
		// Получаем JWT из куков
		Req_ctx := r.Context()
		jwtCookie, err := r.Cookie("jwt_token")
		if err != nil {
			// next.ServeHTTP(w, r)
			// return
			// Если JWT отсутствует, проверяем refresh токен
			refreshCookie, err := r.Cookie("refresh_token")
			if err != nil {
				// http.Error(w, "Unauthorized: No tokens provided", http.StatusUnauthorized)
				// next.ServeHTTP(w, r)
				log.Printf("No JWT or refresh token found")
				return
			}

			// Проверяем refresh токен и получаем новый JWT
			newJWToken, err := CheckRefreshTokenTTL(refreshCookie.Value)
			if err != nil {
				//http.Error(w, "Unauthorized: Invalid refresh token", http.StatusUnauthorized)
				// next.ServeHTTP(w, r)
				log.Printf("Refresh token check failed: %v", err)
				http.Error(w, "Session expired. Please login again.", http.StatusUnauthorized)
				return
			}

			// Устанавливаем новый JWT в куки
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt_token",
				Value:    newJWToken,
				Expires:  time.Now().Add(TokenTTL),
				HttpOnly: true,
				//Secure:   true,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
			})

			// Продолжаем выполнение запроса с новым JWT
			// r.Header.Set("Authorization", "Bearer "+newJWToken)
			userID, err := ParseToken(newJWToken)
			ctx := context.WithValue(Req_ctx, "userID", userID)

			if err != nil || userID == 0 {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Если JWT есть, проверяем его валидность
		userID, err := ParseToken(jwtCookie.Value)
		if err != nil || userID == 0 {
			// Если JWT невалиден, проверяем refresh токен
			refreshCookie, err := r.Cookie("refresh_token")
			if err != nil {
				log.Printf("Invalid JWT and no refresh token found")
				http.Error(w, "Unauthorized: Invalid JWT and no refresh token", http.StatusUnauthorized)
				return
			}

			// Проверяем refresh токен и получаем новый JWT
			newJWToken, err := CheckRefreshTokenTTL(refreshCookie.Value)
			if err != nil {
				log.Printf("Refresh token check failed: %v", err)
				http.Error(w, "Unauthorized: Invalid refresh token", http.StatusUnauthorized)
				return
			}

			// Устанавливаем новый JWT в куки
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt_token",
				Value:    newJWToken,
				Expires:  time.Now().Add(TokenTTL),
				HttpOnly: true,
				//Secure:   true,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
			})

			// // Продолжаем выполнение запроса с новым JWT
			//r.Header.Set("Authorization", "Bearer "+newJWToken)
			userID, err := ParseToken(newJWToken)
			if err != nil || userID == 0 {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(Req_ctx, "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Если JWT валиден, продолжаем выполнение запроса
		ctx := context.WithValue(Req_ctx, "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

}

func CheckRefreshTokenTTL(refresh_token string) (string, error) {
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

	log.Printf("Checking refresh token: %s", refresh_token)

	var id int
	var user_id int
	var expires_at time.Time
	query := "SELECT id, user_id, expires_at FROM refresh_tokens WHERE token = $1"
	err = db.QueryRow(ctx, query, refresh_token).Scan(&id, &user_id, &expires_at)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Refresh token not found: %s", refresh_token)
			return "", errors.New("invalid refresh token")
		}
		log.Printf("DB query error: %v", err)
		return "", fmt.Errorf("internal server error")
	}
	// Проверяем срок действия refresh токена
	if time.Now().After(expires_at) {
		_, err = db.Exec(ctx, "DELETE FROM refresh_tokens WHERE id = $1", id)
		if err != nil {
			return "", fmt.Errorf("failed to delete expired token: %v", err)
		}
		// Начинаем транзакцию
		tx, err := db.Begin(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to start transaction: %v", err)
		}
		defer tx.Rollback(ctx)
		newRefreshToken, err := GenerateRefreshToken()
		if err != nil {
			return "", fmt.Errorf("failed to generate refresh token: %v", err)
		}
		newExpiresAt := time.Now().Add(RefreshTokenTTL)

		// Сохраняем новый refresh токен
		_, err = tx.Exec(ctx,
			"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
			user_id, newRefreshToken, newExpiresAt,
		)
		if err != nil {
			log.Printf("Error saving new refresh token: %v", err)
			return "", fmt.Errorf("failed to save new refresh token")
		}

		// Коммитим транзакцию
		if err := tx.Commit(ctx); err != nil {
			return "", fmt.Errorf("transaction commit failed: %v", err)
		}

		log.Printf("Successfully refreshed tokens for user %d", user_id)
	}

	// Получаем email пользователя
	var email string
	var password string
	err = db.QueryRow(ctx, "SELECT email FROM users WHERE id = $1", user_id).Scan(&email)
	if err != nil {
		return "", fmt.Errorf("failed to get user email: %v", err)
	}
	err = db.QueryRow(ctx, "SELECT password FROM users WHERE id = $1", user_id).Scan(&password)
	if err != nil {
		return "", fmt.Errorf("failed to get user password: %v", err)
	}

	// Начинаем транзакцию
	// tx, err := db.Begin(ctx)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to start transaction: %v", err)
	// }
	// defer tx.Rollback(ctx)

	// // Удаляем старый токен в рамках транзакции
	// _, err = tx.Exec(ctx, "DELETE FROM refresh_tokens WHERE id = $1", id)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to delete old token: %v", err)
	// }

	// Генерируем новые токены
	newJWToken, err := GenerateJWToken(email, password)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %v", err)
	}

	return newJWToken, nil
}
