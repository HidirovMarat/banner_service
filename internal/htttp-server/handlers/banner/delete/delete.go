package delete

import (
	"errors"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"banner-service/internal/lib/logger/sl"
	"banner-service/internal/storage"
	"banner-service/internal/storage/post"
	"context"
)

// BannerDeleter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ContentGetter
type BannerDeleter interface {
	DeletBanner(ctx context.Context, id int64) error
}

func New(ctx context.Context, log *slog.Logger, bannerDeleter BannerDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.banner.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		accessLevel := r.Context().Value("accessLever")

		if accessLevel != post.Admin {
			http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
			return
		}

		sId := chi.URLParam(r, "id")

		if sId == "" {
			log.Info("id empty")
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(sId, 10, 64)

		if err != nil {
			log.Info("id is not number")
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		err = bannerDeleter.DeletBanner(ctx, id)

		if errors.Is(err, storage.ErrBannerNotFound) {
			log.Error("Not id, delete banner", sl.Err(err))
			http.Error(w, storage.ErrBannerNotFound.Error(), http.StatusNotFound)
			return
		}

		if err != nil {
			log.Error("failed to delete banners", sl.Err(err))
			http.Error(w, storage.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Баннер успешно удален"))
	}
}
