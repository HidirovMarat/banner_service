package get

import (
	"errors"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"banner-service/internal/entity"
	resp "banner-service/internal/lib/api/response"
	"banner-service/internal/lib/logger/sl"
	"banner-service/internal/storage"
	"context"
)

// ContentGetter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ContentGetter
type ContentGetter interface {
	GetContent(ctx context.Context, feature_id int64, tag_id int64) (entity.Content, error)
}

type Response struct {
	Content entity.Content `json:"content,omitempty"`
}

func New(ctx context.Context, log *slog.Logger, contentGetter ContentGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.content.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		sFeat_id := r.URL.Query().Get("feature_id")
		sTag_id := r.URL.Query().Get("tag_id")

		if sFeat_id == "" || sTag_id == "" {
			log.Info("feature_id or tag_id is empty")
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		feature_id, err := strconv.ParseInt(sFeat_id, 10, 64)

		if err != nil {
			log.Info("feature_id not number")
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		tag_id, err := strconv.ParseInt(sTag_id, 10, 64)

		if err != nil {
			log.Info("tag_id not number")
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		resContent, err := contentGetter.GetContent(ctx, feature_id, tag_id)

		if errors.Is(err, storage.ErrContentNotFound) {
			log.Info("content not found", sl.Err(err))
			http.Error(w, storage.ErrContentNotFound.Error(), http.StatusNotFound)
			return
		}

		if err != nil {
			log.Error("failed to get content", sl.Err(err))

			render.JSON(w, r, resp.Error(storage.ErrInternalServer.Error()))
			http.Error(w, storage.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, resContent)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, content entity.Content) {
	render.JSON(w, r, Response{
		Content: content,
	})
}
