package service_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kristina71/bitlytest/mocks"
	"github.com/kristina71/bitlytest/pkg/models"
	"github.com/kristina71/bitlytest/pkg/service"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name         string
	expectedUrls []models.Url
	expectedUrl  models.Url
	wantErr      bool
	prepare      func(t *testing.T) string
}

func TestGetBySmallUrl(t *testing.T) {
	testCases := []testCase{
		{
			name: "Get link by small url",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "dfgdfg",
				OriginUrl: "http://google.com",
			},
		},
		{
			name: "Get link by empty small url",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "",
				OriginUrl: "http://google.com",
			},
		},
		/*{
			name: "Get link by small url with slash",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "/dfgdfg",
				OriginUrl: "http://google.com",
			},
		},*/
		/*{
			name: "Get link with slash",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "/",
				OriginUrl: "http://google.com",
			},
		},*/
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			storage := &mocks.Storage{}
			service := service.New(storage)

			storage.On("GetBySmallUrl", models.Url{
				SmallUrl: testCase.expectedUrl.SmallUrl,
			}).Return(testCase.expectedUrl, nil)

			req := httptest.NewRequest(http.MethodGet, "http://localhost:8000/"+testCase.expectedUrl.SmallUrl, nil)
			recorder := httptest.NewRecorder()

			service.GetUrl(recorder, req)
			resp := recorder.Result()
			require.NotNil(t, resp)

			require.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
			require.Equal(t, testCase.expectedUrl.OriginUrl, resp.Header.Get("Location"))
		})
	}
}

func TestGetUrl(t *testing.T) {

	testCases := []testCase{
		{
			name: "Get urls with several models",
			expectedUrls: []models.Url{
				{
					Id:        1,
					SmallUrl:  "dfgdfg",
					OriginUrl: "http://google.com",
				},
				{
					Id:        2,
					SmallUrl:  "dfgddsfdsffg",
					OriginUrl: "http://yandex.ru",
				},
			},
		},
		{
			name: "Get urls with one model",
			expectedUrls: []models.Url{
				{
					Id:        1,
					SmallUrl:  "dfgdfg",
					OriginUrl: "http://google.com",
				},
			},
		},
		{
			name: "Get urls with small url",
			expectedUrls: []models.Url{
				{
					Id:        1,
					SmallUrl:  "",
					OriginUrl: "http://google.com",
				},
			},
		},
		{
			name: "Get urls with empty origin url",
			expectedUrls: []models.Url{
				{
					Id:        1,
					SmallUrl:  "fgdfgdfgfdg",
					OriginUrl: "",
				},
			},
		},
		{
			name:         "Get urls with empty model",
			expectedUrls: []models.Url{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			storage := &mocks.Storage{}
			service := service.New(storage)
			storage.On("Get").Return(testCase.expectedUrls, nil)

			var jsonStr = []byte(`{}`)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8000/all", bytes.NewBuffer(jsonStr))
			recorder := httptest.NewRecorder()

			service.GetAllUrl(recorder, req)
			resp := recorder.Result()
			require.NotNil(t, resp)

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			expectedBody, err := json.Marshal(testCase.expectedUrls)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, string(expectedBody), string(body))
		})
	}
}

func TestUpdateUrl(t *testing.T) {
	testCases := []testCase{
		{
			name: "Update url",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "dfgdfg",
				OriginUrl: "http://google.com",
			},
		},
		{
			name: "Update url with empty small url",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "",
				OriginUrl: "http://google.com",
			},
		},
		{
			name: "Update url with empty small url with slash",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "/",
				OriginUrl: "http://google.com",
			},
		},
		{
			name: "Update url with incorrect url",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "/",
				OriginUrl: "dfdfdf",
			},
		},
		{
			name: "Update url with empty origin url",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "fdggdfgfd",
				OriginUrl: "",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			storage := &mocks.Storage{}
			service := service.New(storage)

			storage.On("Update", testCase.expectedUrl).Return(nil)

			b, err := json.Marshal(models.Url{Id: testCase.expectedUrl.Id,
				SmallUrl: testCase.expectedUrl.SmallUrl, OriginUrl: testCase.expectedUrl.OriginUrl})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8000/edit", bytes.NewBuffer(b))
			recorder := httptest.NewRecorder()

			service.UpdateUrl(recorder, req)
			resp := recorder.Result()
			require.NotNil(t, resp)

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			url := models.Url{}
			err = json.Unmarshal(body, &url)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, testCase.expectedUrl, url)
		})
	}
}

func TestInsertUrl(t *testing.T) {
	testCases := []testCase{
		{
			name: "Create url",
			expectedUrl: models.Url{
				Id:        2,
				SmallUrl:  "dfgddsfdsffg",
				OriginUrl: "http://yandex.ru",
			},
		},
		{
			name: "Create url with empty small url",
			expectedUrl: models.Url{
				Id:        2,
				SmallUrl:  "",
				OriginUrl: "http://yandex.ru",
			},
		},
		{
			name: "Create url with empty small url with slash",
			expectedUrl: models.Url{
				Id:        2,
				SmallUrl:  "/",
				OriginUrl: "http://yandex.ru",
			},
		},
		{
			name: "Create url with empty origin url",
			expectedUrl: models.Url{
				Id:        2,
				SmallUrl:  "fgdfgf",
				OriginUrl: "",
			},
		},
		{
			name: "Create url with incorrect origin url",
			expectedUrl: models.Url{
				Id:        2,
				SmallUrl:  "fgdfgf",
				OriginUrl: "gfdgdfgfd",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			storage := &mocks.Storage{}
			service := service.New(storage)

			storage.On("Insert", models.Url{
				SmallUrl: testCase.expectedUrl.SmallUrl, OriginUrl: testCase.expectedUrl.OriginUrl,
			}).Return(testCase.expectedUrl.Id, nil)

			b, err := json.Marshal(models.Url{
				SmallUrl: testCase.expectedUrl.SmallUrl, OriginUrl: testCase.expectedUrl.OriginUrl})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8000/create", bytes.NewBuffer(b))
			recorder := httptest.NewRecorder()

			service.CreateUrl(recorder, req)
			resp := recorder.Result()
			require.NotNil(t, resp)

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			url := models.Url{}
			err = json.Unmarshal(body, &url)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, testCase.expectedUrl, url)
		})
	}
}

func TestDeleteUrl(t *testing.T) {
	testCases := []testCase{
		{
			name: "Delete url",
			expectedUrl: models.Url{
				Id: 2,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			storage := &mocks.Storage{}
			service := service.New(storage)

			storage.On("Delete", models.Url{Id: testCase.expectedUrl.Id}).Return(nil)

			jsonStr, err := json.Marshal(models.Url{Id: testCase.expectedUrl.Id})
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8000/delete", bytes.NewBuffer(jsonStr))
			recorder := httptest.NewRecorder()

			service.DeleteUrl(recorder, req)
			resp := recorder.Result()
			require.NotNil(t, resp)
		})
	}
}
