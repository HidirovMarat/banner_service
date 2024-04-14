package tests

import (
	//	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"

	//	"github.com/stretchr/testify/require"
	"banner-service/internal/entity"

	"banner-service/internal/htttp-server/handlers/banner/create"
	//	"banner-service/internal/lig/api"
	//"banner-service/internal/lib/random"
)

const (
	host = "localhost:8082"
)

func TestBanner_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/banner").
		WithJSON(create.Request{
			Content: entity.Content{
				Title: gofakeit.Book().Title,
				Text:  gofakeit.Book().Genre,
				Url:   gofakeit.URL(),
			},
			Feature_id: int64(gofakeit.Int32()),
			Tag_ids:    []int64{int64(gofakeit.Int32()), int64(gofakeit.Int32()), int64(gofakeit.Int32()), int64(gofakeit.Int32())},
			Is_active:  gofakeit.Bool(),
		}).
		WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxODA5MjQsInN1YiI6ImFkbWluIn0.LqnReEg1smLzxYIuKG1GHzXYjRg572IlGvMO28hKNfU").
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("id")
}

//nolint:funlen
