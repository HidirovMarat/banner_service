package get

import (
	"errors"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi/v5"
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
type BannerGetter interface {
	GetBanners(ctx context.Context, feature_id, tag_id *int64, offset, limit *int) ([]entity.Banner, error)
}

type Response struct {
	resp.Response
	Banners []entity.Banner `json:"content,omitempty"`
}

func New(ctx context.Context, log *slog.Logger, bannerGetter BannerGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.banner.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		sFeat_id := chi.URLParam(r, "featuer_id")
		sTag_id := chi.URLParam(r, "tag_id")
		sOffset := chi.URLParam(r, "offset")
		sLimit := chi.URLParam(r, "limit")

		var feature_id, tag_id *int64
		var offset, limit *int

		if sFeat_id != "" {
			log.Info("feature_id is empty")

			f, _ := strconv.ParseInt(sFeat_id, 10, 64)

			feature_id = &f
		}

		if sTag_id != "" {
			log.Info("tag_id is empty")

			t, _ := strconv.ParseInt(sTag_id, 10, 64)

			tag_id = &t
		}

		if sOffset != "" {
			log.Info("offset is empty")

			o, _ := strconv.Atoi(sOffset)

			offset = &o
		}

		if sLimit != "" {
			log.Info("feature_id is empty")

			l, _ := strconv.Atoi(sLimit)

			limit = &l
		}

		resBanners, err := bannerGetter.GetBanners(ctx, feature_id, tag_id, offset, limit)

		if errors.Is(err, storage.ErrBannerNotFound) {
			log.Error("not id in banners", sl.Err(err))

			render.JSON(w, r, resp.Error(storage.ErrBannerNotFound.Error()))

			return
		}

		if err != nil {
			log.Error("failed to get banners", sl.Err(err))

			render.JSON(w, r, resp.Error(storage.ErrInternalServer.Error()))

			return
		}

		responseOK(w, r, resBanners)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, banners []entity.Banner) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Banners:  banners,
	})
}
