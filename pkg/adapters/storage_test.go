package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/dailymotion/allure-go"
	"github.com/kristina71/bitlytest/pkg/adapters"
	"github.com/kristina71/bitlytest/pkg/models"

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

func TestInsertDB(t *testing.T) {
	allure.Test(t,
		allure.Description("Insert data in DB"),
		allure.Action(func() {
			db, mock, err := sqlxmock.Newx()
			require.NoError(t, err)

			defer db.Close()

			storage := adapters.New(db)

			testCases := []testCase{
				{
					name: "OK",
					url: models.Url{
						SmallUrl:  "xyz",
						OriginUrl: "dsfsdfds",
						CreatedAt: time.Now(),
						UpdateAt:  time.Now(),
					},
					mock: func(tc *testCase) {
						rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
						mock.ExpectQuery("INSERT INTO bitlytest").WithArgs(tc.url.SmallUrl, tc.url.OriginUrl, tc.url.CreatedAt, tc.url.UpdateAt).WillReturnRows(rows)
					},
					id:      1,
					wantErr: false,
				},
				{
					name: "Inser empty fields",
					url: models.Url{
						SmallUrl:  "",
						OriginUrl: "",
						CreatedAt: time.Now(),
						UpdateAt:  time.Now(),
					},
					mock: func(tc *testCase) {
						rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
						mock.ExpectQuery("INSERT INTO bitlytest").WithArgs(tc.url.SmallUrl, tc.url.OriginUrl, tc.url.CreatedAt, tc.url.UpdateAt).WillReturnRows(rows)
					},
					wantErr: true,
				},
			}

			for _, testCase := range testCases {
				t.Run(testCase.name, func(t *testing.T) {
					mockData(testCase)

					allure.Step(allure.Description("Insert data and check result"), allure.Action(func() {
						insertId, err := storage.Insert(context.TODO(), testCase.url)

						require.NoError(t, err)
						if testCase.wantErr != true {
							require.Equal(t, testCase.id, insertId)
						}
					}))
				})
			}
		}))
}

func TestSelectDB(t *testing.T) {
	allure.Test(t,
		allure.Description("Select data in DB"),
		allure.Action(func() {
			db, mock, err := sqlxmock.Newx()
			require.NoError(t, err)

			defer db.Close()

			storage := adapters.New(db)

			testCases := []testCase{
				{
					name: "OK",
					url: models.Url{
						Id:        1,
						SmallUrl:  "xyz",
						OriginUrl: "dsfsdfds",
						CreatedAt: time.Now(),
						UpdateAt:  time.Now(),
					},
					mock: func(tc *testCase) {
						rows := sqlxmock.NewRows([]string{"id", "small_url", "origin_url", "created_at", "updated_at"}).
							AddRow(tc.url.Id, tc.url.SmallUrl, tc.url.OriginUrl, tc.url.CreatedAt, tc.url.UpdateAt)
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
					mockData(testCase)

					allure.Step(allure.Description("Get data by small url and check result"), allure.Action(func() {
						modelUrl, err := storage.GetBySmallUrl(context.TODO(), testCase.url)

						require.NoError(t, err)
						require.Equal(t, testCase.url, modelUrl)
					}))
				})
			}
		}))
}

func TestDeleteDB(t *testing.T) {
	allure.Test(t,
		allure.Description("Delete data in DB"),
		allure.Action(func() {
			db, mock, err := sqlxmock.Newx()
			require.NoError(t, err)

			defer db.Close()

			storage := adapters.New(db)

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
					mockData(testCase)

					allure.Step(allure.Description("Delete data and check result"), allure.Action(func() {
						err := storage.Delete(context.TODO(), testCase.url)

						require.NoError(t, err)
					}))
				})
			}
		}))
}

func TestUpdateDB(t *testing.T) {
	allure.Test(t,
		allure.Description("Delete data in DB"),
		allure.Action(func() {
			db, mock, err := sqlxmock.Newx()
			require.NoError(t, err)

			defer db.Close()

			storage := adapters.New(db)

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
					mockData(testCase)

					allure.Step(allure.Description("Update data and check result"), allure.Action(func() {
						err := storage.Update(context.TODO(), testCase.url)
						require.NoError(t, err)
					}))
				})
			}
		}))
}

func mockData(testCase testCase) {
	allure.Step(allure.Description("Mock data"), allure.Action(func() {
		testCase.mock(&testCase)
	}))
}
