package create_test

import (
	"bytes"
	"context"

 	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"banner-service/internal/entity"
	"banner-service/internal/http-server/handlers/banner/create"
	"banner-service/internal/http-server/handlers/banner/create/mocks"
	"banner-service/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name        string
		accessLevel string
		feature_id  int64
		tag_ids     []int64
		is_active   bool
		content     entity.Content
		respError   string
		mockError   error
		statucCode  int
	}{
		{
			name:        "1 Test",
			accessLevel: "admin",
			feature_id:  3,
			tag_ids:     []int64{2, 3, 5, 23},
			is_active:   true,
			content: entity.Content{
				Title: "title",
				Text:  "Text",
				Url:   "ds",
			},
			statucCode: 200,
		},
		{
			name:       "2 Test",
			feature_id: 4,
			tag_ids:    []int64{2, 3, 23},
			is_active:  true,
			content: entity.Content{
				Title: "title",
				Text:  "Text",
				Url:   "ds",
			},
			respError:  "Пользователь не имеет доступа",
			statucCode: 403,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bannnerCreateMock := mocks.NewBannerCreate(t)

			if tc.respError == "" || tc.mockError != nil {
				bannnerCreateMock.On("CreateBanner", context.Background(), tc.feature_id, tc.tag_ids, tc.is_active, tc.content).
					Return(int64(3), tc.mockError).
					Once()
			}

			handler := create.New(context.Background(), slogdiscard.NewDiscardLogger(), bannnerCreateMock)

			sTags := "["

			for index, val := range tc.tag_ids {

				sTags += strconv.FormatInt(val, 10)

				if index != len(tc.tag_ids)-1 {
					sTags += ","
				}
			}

			sTags += "]"

			input := fmt.Sprintf(
				`{
				"feature_id": %d,
			 	"tag_ids": %s,
			  	"is_active": %t,
			   	"content":  {
					"text": "%s",
					"title":"%s",
					"url":"%s"
				}
			}`,
				tc.feature_id, sTags, tc.is_active, tc.content.Text, tc.content.Title, tc.content.Url)

			req, err := http.NewRequest(http.MethodPost, "/banner", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			req = req.WithContext(context.WithValue(req.Context(), "accessLever", tc.accessLevel))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.statucCode)
		})
	}
}
