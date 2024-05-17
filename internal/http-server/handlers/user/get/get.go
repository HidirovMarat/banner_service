package get

import (
	"errors"
	"net/http"
	"time"

	"log/slog"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"

	"banner-service/internal/lib/logger/sl"
	"banner-service/internal/storage"
	"context"
)

// ContentGetter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ContentGetter
type UserGetter interface {
	GetAccessLevel(ctx context.Context, login string, password string) (string, error)
}

type Response struct {
	JWT string `json:"content,omitempty"`
}

func New(ctx context.Context, log *slog.Logger, userGetter UserGetter, signingKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")

		if login == "" || password == "" {
			log.Info("feature_id or tag_id is empty")
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		res, err := userGetter.GetAccessLevel(ctx, login, password)

		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("not user set login, password", sl.Err(err))
			http.Error(w, storage.ErrUserNotFound.Error(), http.StatusInternalServerError)
			return
		}

		if err != nil {
			log.Error("failed to get user", sl.Err(err))
			http.Error(w, storage.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

		if res == "" {
			log.Info("found not user", login, password)
			http.Error(w, storage.ErrUserNotFound.Error(), http.StatusNotFound)
			return
		}

		token, err := GetToken(res, signingKey)

		if err != nil {
			log.Info("make jwt-token error", login, password)
			http.Error(w, storage.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", "Bearer "+token)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Токен успешно отправлен в заголовке"))
	}
}

func GetToken(accessLevel string, signingKey string) (string, error) {
	// Создаем структуру для хранения стандартных параметров токена
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Время истечения токена через 24 часа
	}
	claims.Subject = accessLevel
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
