package create

import (
	"banner-service/internal/entity"
	"banner-service/internal/storage"
	"context"
	"errors"
	"io"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "banner-service/internal/lib/api/response"
	"banner-service/internal/lib/logger/sl"
	"banner-service/internal/storage/post"
)

type Request struct {
	Content    entity.Content `json:"content" validate:"required"`
	Feature_id int64          `json:"feature_id" validate:"required"`
	Is_active  bool           `json:"is_active" validate:"required"`
	Tag_ids    []int64        `json:"tag_ids" validate:"required"`
}

type Response struct {
	Id int64 `json:"id,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=BannerCreate
type BannerCreate interface {
	CreateBanner(ctx context.Context, feature_id int64, tag_ids []int64, is_active bool, content entity.Content) (int64, error)
}

func New(context context.Context, log *slog.Logger, bannerCreate BannerCreate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.banner.create.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		accessLevel := r.Context().Value("accessLever")

		if accessLevel != post.Admin {
			http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
			return 
		}

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			http.Error(w, storage.ErrIncorrectData.Error(), http.StatusBadRequest)
			return
		}

		id, err := bannerCreate.CreateBanner(context, req.Feature_id, req.Tag_ids, req.Is_active, req.Content)

		if errors.Is(err, storage.ErrIncorrectData) {
			log.Info("banner already exists")
			render.JSON(w, r, resp.Error(storage.ErrIncorrectData.Error()))
			return
		}

		if err != nil {
			log.Error("failed to add banner", sl.Err(err))
			http.Error(w, storage.ErrInternalServer.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("banner added", slog.Int64("id", id))

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int64) {
	render.JSON(w, r, Response{
		Id: id,
	})
}
