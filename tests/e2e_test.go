package tests

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kristina71/bitlytest/pkg/adapters"
	"github.com/kristina71/bitlytest/pkg/config"
	"github.com/kristina71/bitlytest/pkg/endpoints"
	"github.com/kristina71/bitlytest/pkg/models"
	"github.com/kristina71/bitlytest/pkg/repositories"
	"github.com/kristina71/bitlytest/pkg/service"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name             string
	expected_status  int
	body             models.Url
	error_checker    func(t require.TestingT, err error, msgAndArgs ...interface{})
	db_error_checker func(t require.TestingT, err error, msgAndArgs ...interface{})
	insert           bool
}

type testCaseAll struct {
	name             string
	expected_status  int
	body             []models.Url
	error_checker    func(t require.TestingT, err error, msgAndArgs ...interface{})
	db_error_checker func(t require.TestingT, err error, msgAndArgs ...interface{})
}

func TestCreate1(t *testing.T) {
	cfg := config.New()
	db := adapters.DBConnect(cfg)
	adapters := adapters.New(db)
	repo := repositories.New(adapters)
	service := service.New(repo)

	ts := httptest.NewServer(endpoints.New(service))
	defer ts.Close()

	testCases := []testCase{
		{
			name: "Item",
			body: models.Url{
				SmallUrl:  "test22211",
				OriginUrl: "http://google.ru/test22211",
				CreatedAt: time.Now(),
				UpdateAt:  time.Now(),
			},
			expected_status:  http.StatusOK,
			error_checker:    require.NoError,
			db_error_checker: require.NoError,
		},
		{
			name: "Empty small url",
			body: models.Url{
				SmallUrl:  "",
				OriginUrl: "http://google.ru/test",
				CreatedAt: time.Now(),
				UpdateAt:  time.Now(),
			},
			expected_status:  http.StatusBadRequest,
			error_checker:    require.NoError,
			db_error_checker: require.Error,
		},
		{
			name: "Empty origin url",
			body: models.Url{
				SmallUrl:  "fdfdfg",
				OriginUrl: "",
				CreatedAt: time.Now(),
				UpdateAt:  time.Now(),
			},
			expected_status:  http.StatusBadRequest,
			error_checker:    require.NoError,
			db_error_checker: require.Error,
		},
		{
			name: "Incorrect origin url",
			body: models.Url{
				SmallUrl:  "fdfdfg",
				OriginUrl: "/",
				CreatedAt: time.Now(),
				UpdateAt:  time.Now(),
			},
			expected_status:  http.StatusBadRequest,
			error_checker:    require.NoError,
			db_error_checker: require.Error,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				resp, err := CreateItem(ts, testCase.body)
				require.Equal(t, testCase.expected_status, resp.StatusCode)
				testCase.error_checker(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				//не через адаптер а через сервис
				dbResult, err := adapters.GetBySmallUrl(context.Background(), testCase.body)
				testCase.db_error_checker(t, err)

				if err == nil {
					url := models.Url{}
					err = json.Unmarshal(body, &url)
					require.NoError(t, err)

					require.Equal(t, dbResult, url)

					err = DeleteItem(ts, url)
					require.NoError(t, err)
				}
			})
	}
}

func TestGetAll1(t *testing.T) {
	cfg := config.New()
	db := adapters.DBConnect(cfg)
	adapters := adapters.New(db)
	repo := repositories.New(adapters)
	service := service.New(repo)

	ts := httptest.NewServer(endpoints.New(service))
	defer ts.Close()

	testCases := []testCaseAll{
		{
			name: "Edit item",
			body: []models.Url{
				{
					Id:        1,
					SmallUrl:  "dfgdfg",
					OriginUrl: "http://google.com",
					CreatedAt: time.Now(),
					UpdateAt:  time.Now(),
				},
				{
					Id:        2,
					SmallUrl:  "dfgddsfdsffg",
					OriginUrl: "http://yandex.ru",
					CreatedAt: time.Now(),
					UpdateAt:  time.Now(),
				},
			},
			expected_status:  http.StatusOK,
			error_checker:    require.NoError,
			db_error_checker: require.NoError,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {

				for _, bodyItem := range testCase.body {
					_, err := CreateItem(ts, bodyItem)
					require.NoError(t, err)
				}

				resp, err := ts.Client().Get(ts.URL + "/all")
				require.Equal(t, testCase.expected_status, resp.StatusCode)
				testCase.error_checker(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				//не через адаптер а через сервис
				dbResult, err := adapters.Get(context.TODO())
				testCase.db_error_checker(t, err)

				if err == nil {
					url := []models.Url{}
					err = json.Unmarshal(body, &url)
					require.NoError(t, err)

					require.Equal(t, dbResult, url)

					err = DeleteItem(ts, url[0])
					require.NoError(t, err)
				}
			})
	}
}
