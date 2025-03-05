package service

import (
	"context"
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из куки
		Req_ctx := r.Context()
		cookie, err := r.Cookie("jwt_token")
		if err != nil {
			//http.Redirect(w, r, "/sign_in", http.StatusUnauthorized)
			next.ServeHTTP(w, r)
			return
		}

		// Парсим токен
		userID, err := ParseToken(cookie.Value)
		if err != nil || userID == 0 {
			//http.Redirect(w, r, "/sign_in", http.StatusUnauthorized)
			next.ServeHTTP(w, r)
			return
		}

		// Добавляем userID в контекст запроса
		ctx := context.WithValue(Req_ctx, "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
