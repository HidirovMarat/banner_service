package jwt

import (
	"context"
	"net/http"
	"strings"

	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	AccessLevel string `json:"accessLevel"`
	jwt.StandardClaims
}

func New(signingKey string) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				// Если заголовок не установлен или не начинается с "Bearer ", возвращаем ошибку
				http.Error(w, "Отсутствует или неверный заголовок Authorization", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Парсинг токена с учетом структуры CustomClaims
			token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(signingKey), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Отсутствует или неверный token"))
				return
			}
			var ctx context.Context
			if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid && (time.Now().Before(time.Unix(claims.ExpiresAt, 0))) {
				ctx = context.WithValue(r.Context(), "accessLever", claims.Subject)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Пользователь не авторизован"))
			}
		}

		return http.HandlerFunc(fn)
	}
}
