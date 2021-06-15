package storage_test

import (
	"bitlytest/pkg/models"
	"bitlytest/pkg/storage"
	"testing"

	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

type testCase struct {
	name    string
	url     models.Url
	mock    func(tc *testCase)
	id      uint16
	wantErr bool
}

//дописать негативные проверки на err
func TestInsertDB(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	require.NoError(t, err)

	defer db.Close()

	storage := storage.New(db)

	testCases := []testCase{
		{
			name: "OK",
			url: models.Url{
				SmallUrl:  "xyz",
				OriginUrl: "dsfsdfds",
			},
			mock: func(tc *testCase) {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO bitlytest").WithArgs(tc.url.SmallUrl, tc.url.OriginUrl).WillReturnRows(rows)
			},
			id:      1,
			wantErr: false,
		},
		{
			name: "Inser empty fields",
			url: models.Url{
				SmallUrl:  "",
				OriginUrl: "",
			},
			mock: func(tc *testCase) {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO bitlytest").WithArgs(tc.url.SmallUrl, tc.url.OriginUrl).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock(&testCase)
			insertId, err := storage.Insert(testCase.url)

			require.NoError(t, err)
			if testCase.wantErr != true {
				require.Equal(t, testCase.id, insertId)
			}
		})
	}
}

func TestSelectDB(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	require.NoError(t, err)

	defer db.Close()

	storage := storage.New(db)

	testCases := []testCase{
		{
			name: "OK",
			url: models.Url{
				Id:        1,
				SmallUrl:  "xyz",
				OriginUrl: "dsfsdfds",
			},
			mock: func(tc *testCase) {
				rows := sqlxmock.NewRows([]string{"id", "small_url", "origin_url"}).
					AddRow(tc.url.Id, tc.url.SmallUrl, tc.url.OriginUrl)
				mock.ExpectQuery("^SELECT (.+) FROM bitlytest WHERE small_url = \\$1").
					WithArgs(tc.url.SmallUrl).
					WillReturnRows(rows)
			},
			id:      1,
			wantErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock(&testCase)

			modelUrl, err := storage.GetBySmallUrl(testCase.url)

			require.NoError(t, err)
			require.Equal(t, testCase.url, modelUrl)
		})
	}
}
func TestDeleteDB(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	require.NoError(t, err)

	defer db.Close()

	storage := storage.New(db)

	testCases := []testCase{
		{
			name: "OK",
			url: models.Url{
				Id: 1,
			},
			mock: func(tc *testCase) {
				mock.ExpectExec("^DELETE FROM bitlytest WHERE id = \\$1").
					WithArgs(tc.url.Id).WillReturnResult(sqlxmock.NewResult(1, 1))

			},
			id:      1,
			wantErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock(&testCase)

			err := storage.Delete(testCase.url)

			require.NoError(t, err)
		})
	}
}

func TestUpdateDB(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	require.NoError(t, err)

	defer db.Close()

	storage := storage.New(db)

	testCases := []testCase{
		{
			name: "OK",
			url: models.Url{
				Id:        1,
				SmallUrl:  "xyz",
				OriginUrl: "dsfsdfds",
			},
			mock: func(tc *testCase) {
				mock.ExpectExec("^UPDATE bitlytest SET small_url = \\$1, origin_url = \\$2 WHERE id = \\$3").
					WithArgs(tc.url.SmallUrl,
						tc.url.OriginUrl,
						tc.url.Id,
					).WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			id:      1,
			wantErr: false,
		},
		{
			name: "Update with empty fields",
			url: models.Url{
				Id:        1,
				SmallUrl:  "",
				OriginUrl: "",
			},
			mock: func(tc *testCase) {
				mock.ExpectExec("^UPDATE bitlytest SET small_url = \\$1, origin_url = \\$2 WHERE id = \\$3").
					WithArgs(tc.url.SmallUrl,
						tc.url.OriginUrl,
						tc.url.Id,
					).WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			id:      1,
			wantErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock(&testCase)

			err := storage.Update(testCase.url)
			require.NoError(t, err)
		})
	}
}
