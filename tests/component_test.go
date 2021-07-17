package tests

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/kristina71/bitlytest/pkg/adapters"
	"github.com/kristina71/bitlytest/pkg/config"
	"github.com/kristina71/bitlytest/pkg/endpoints"
	"github.com/kristina71/bitlytest/pkg/models"
	"github.com/kristina71/bitlytest/pkg/repositories"
	"github.com/kristina71/bitlytest/pkg/service"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
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
				OriginUrl: "http://google.ru",
			},
			expected_status:  http.StatusOK,
			error_checker:    require.NoError,
			db_error_checker: require.NoError,
		},
		{
			name: "Not found origin url",
			body: models.Url{
				SmallUrl:  "test22211",
				OriginUrl: "http://google.ru/gfdgdfgdf",
			},
			expected_status:  http.StatusBadRequest,
			error_checker:    require.NoError,
			db_error_checker: require.Error,
		},
		{
			name: "Empty small url",
			body: models.Url{
				SmallUrl:  "",
				OriginUrl: "http://google.ru/test",
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
			},
			expected_status:  http.StatusBadRequest,
			error_checker:    require.NoError,
			db_error_checker: require.Error,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				defer clean(db, t)

				resp, err := CreateItem(ts, testCase.body)
				require.Equal(t, testCase.expected_status, resp.StatusCode)
				testCase.error_checker(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				dbResult, err := adapters.GetBySmallUrl(context.Background(), testCase.body)
				testCase.db_error_checker(t, err)

				if err == nil {
					url := models.Url{}
					err = json.Unmarshal(body, &url)
					require.NoError(t, err)

					require.Equal(t, dbResult, url)
				}
			})
	}
}

func TestGetBySmallUrl(t *testing.T) {
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
				OriginUrl: "http://google.ru",
			},
			expected_status:  http.StatusPermanentRedirect,
			error_checker:    require.NoError,
			db_error_checker: require.NoError,
			insert:           true,
		},
		{
			name: "Not found",
			body: models.Url{
				SmallUrl:  "test22211",
				OriginUrl: "http://google.ru",
			},
			expected_status: http.StatusNotFound,
			insert:          false,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				defer clean(db, t)

				if testCase.insert == true {
					_, err := adapters.Insert(context.Background(), testCase.body)
					require.NoError(t, err)
				}

				ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}

				resp, err := ts.Client().Get(ts.URL + "/" + testCase.body.SmallUrl)
				require.NoError(t, err)

				defer resp.Body.Close()
				require.Equal(t, testCase.expected_status, resp.StatusCode)

				if testCase.insert == true {
					require.Equal(t, testCase.body.OriginUrl, resp.Header.Get("location"))
				}
			})
	}
}

func TestUpdate(t *testing.T) {
	cfg := config.New()
	db := adapters.DBConnect(cfg)
	adapters := adapters.New(db)
	repo := repositories.New(adapters)
	service := service.New(repo)

	ts := httptest.NewServer(endpoints.New(service))
	defer ts.Close()

	testCases := []testCase{
		{
			name: "Update item",
			body: models.Url{
				SmallUrl:  "test22211",
				OriginUrl: "http://google.ru",
			},
			expected_status:  http.StatusPermanentRedirect,
			error_checker:    require.NoError,
			db_error_checker: require.NoError,
			insert:           true,
		},
		{
			name: "Update item with incorrect id",
			body: models.Url{
				Id:        0,
				SmallUrl:  "test22211",
				OriginUrl: "http://google.ru",
			},
			expected_status:  http.StatusPermanentRedirect,
			error_checker:    require.NoError,
			db_error_checker: require.NoError,
			insert:           false,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				defer clean(db, t)

				if testCase.insert == true {
					id, err := adapters.Insert(context.Background(), testCase.body)
					testCase.db_error_checker(t, err)
					testCase.body.Id = id
				}

				testCase.body.SmallUrl = testCase.body.SmallUrl + "11"
				resp, err := EditItem(ts, testCase.body)
				require.NoError(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				url := models.Url{}
				err = json.Unmarshal(body, &url)
				require.NoError(t, err)

				item, err := adapters.GetBySmallUrl(context.Background(), testCase.body)
				if testCase.insert == true {
					require.NoError(t, err)

					require.Equal(t, item.SmallUrl, url.SmallUrl)
					require.Equal(t, item.OriginUrl, url.OriginUrl)
					require.Equal(t, item.Id, url.Id)

					require.Equal(t, testCase.body.SmallUrl, url.SmallUrl)
					require.Equal(t, testCase.body.OriginUrl, url.OriginUrl)
					require.Equal(t, testCase.body.Id, url.Id)
				} else {
					require.Error(t, err)
				}
			})
	}
}

func TestDelete(t *testing.T) {
	cfg := config.New()
	db := adapters.DBConnect(cfg)
	adapters := adapters.New(db)
	repo := repositories.New(adapters)
	service := service.New(repo)

	ts := httptest.NewServer(endpoints.New(service))
	defer ts.Close()

	testCases := []testCase{
		{
			name: "Delete item",
			body: models.Url{
				SmallUrl:  "test22211",
				OriginUrl: "http://google.ru",
			},
			insert:           true,
			expected_status:  http.StatusPermanentRedirect,
			db_error_checker: require.NoError,
		},
		{
			name: "Try to delete not existing url",
			body: models.Url{
				SmallUrl:  "test22211",
				OriginUrl: "http://google.ru",
			},
			insert:          false,
			expected_status: http.StatusPermanentRedirect,
		},
		{
			name:            "Try to delete with empty body",
			body:            models.Url{},
			insert:          false,
			expected_status: http.StatusPermanentRedirect,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				defer clean(db, t)

				if testCase.insert == true {
					id, err := adapters.Insert(context.Background(), testCase.body)
					testCase.db_error_checker(t, err)
					testCase.body.Id = id
				}

				err := DeleteItem(ts, testCase.body)
				require.NoError(t, err)

				_, err = adapters.GetBySmallUrl(context.Background(), testCase.body)
				require.Error(t, err)
			})
	}
}

func clean(db *sqlx.DB, t *testing.T) {
	_, err := db.Exec("DELETE FROM bitlytest")
	require.NoError(t, err)
}
