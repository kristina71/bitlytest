package tests

import (
	"bitlytest/pkg/config"
	"bitlytest/pkg/models"
	"bitlytest/pkg/router"
	"bitlytest/pkg/storage"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name             string
	expected_status  int
	body             models.Url
	error_checker    func(t require.TestingT, err error, msgAndArgs ...interface{})
	db_error_checker func(t require.TestingT, err error, msgAndArgs ...interface{})
}

type testCaseAll struct {
	name             string
	expected_status  int
	body             []models.Url
	error_checker    func(t require.TestingT, err error, msgAndArgs ...interface{})
	db_error_checker func(t require.TestingT, err error, msgAndArgs ...interface{})
}

func TestCreate(t *testing.T) {
	cfg := config.New()
	db := storage.DBConnect(cfg)
	s := storage.New(db)

	ts := httptest.NewServer(router.NewRouter(db))
	defer ts.Close()

	testCases := []testCase{
		{
			name: "Edit item",
			body: models.Url{
				SmallUrl:  "test2221",
				OriginUrl: "http://google.ru/test2221",
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
				resp, err := CreateItem(ts, testCase.body)
				require.Equal(t, testCase.expected_status, resp.StatusCode)
				testCase.error_checker(t, err)

				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				defer resp.Body.Close()

				dbResult, err := s.GetBySmallUrl(testCase.body)
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

func TestGetAll(t *testing.T) {
	cfg := config.New()
	db := storage.DBConnect(cfg)
	s := storage.New(db)

	ts := httptest.NewServer(router.NewRouter(db))
	defer ts.Close()

	testCases := []testCaseAll{
		{
			name: "Edit item",
			body: []models.Url{
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

				dbResult, err := s.Get()
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
