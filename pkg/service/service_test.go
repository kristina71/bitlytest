package service_test

import (
	"context"
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
		{
			name: "Get link by small url with slash",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "/dfgdfg",
				OriginUrl: "http://google.com",
			},
		},
		{
			name: "Get link with slash",
			expectedUrl: models.Url{
				Id:        1,
				SmallUrl:  "/",
				OriginUrl: "http://google.com",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := &mocks.Repository{}
			service := service.New(repo)

			url := models.Url{
				SmallUrl: testCase.expectedUrl.SmallUrl,
			}
			repo.On("GetBySmallUrl", context.Background(), url).Return(testCase.expectedUrl, nil)

			resUrl, err := service.GetUrl(context.Background(), url)
			require.NoError(t, err)

			require.Equal(t, testCase.expectedUrl, resUrl)
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
			repo := &mocks.Repository{}
			service := service.New(repo)

			repo.On("Get", context.Background()).Return(testCase.expectedUrls, nil)

			resUrls, err := service.GetAllUrl(context.Background())
			require.NoError(t, err)

			require.Equal(t, testCase.expectedUrls, resUrls)
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
			repo := &mocks.Repository{}
			service := service.New(repo)
			expected := testCase.expectedUrl

			repo.On("ValidateUrl", context.Background(), testCase.expectedUrl.OriginUrl).Return(true)
			if testCase.expectedUrl.SmallUrl == "" || testCase.expectedUrl.SmallUrl == "/" {
				repo.On("GenerateUrl", context.Background()).Return("fdfdfdh")
				expected.SmallUrl = "fdfdfdh"
			}
			repo.On("Update", context.Background(), expected).Return(nil)

			resUrl, err := service.UpdateUrl(context.Background(), testCase.expectedUrl)
			require.NoError(t, err)

			require.Equal(t, expected, resUrl)
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
			repo := &mocks.Repository{}
			service := service.New(repo)

			expected := testCase.expectedUrl

			repo.On("ValidateUrl", context.Background(), testCase.expectedUrl.OriginUrl).Return(true)
			if testCase.expectedUrl.SmallUrl == "" || testCase.expectedUrl.SmallUrl == "/" {
				repo.On("GenerateUrl", context.Background()).Return("fdfdfdh")
				expected.SmallUrl = "fdfdfdh"
			}
			repo.On("Insert", context.Background(), expected).Return(testCase.expectedUrl.Id, nil)

			resUrl, err := service.CreateUrl(context.Background(), testCase.expectedUrl)
			require.NoError(t, err)

			require.Equal(t, expected, resUrl)
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
			repo := &mocks.Repository{}
			service := service.New(repo)

			repo.On("Delete", context.Background(), models.Url{Id: testCase.expectedUrl.Id}).Return(nil)

			err := service.DeleteUrl(context.Background(), testCase.expectedUrl)
			require.NoError(t, err)
		})
	}
}
