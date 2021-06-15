package storage

import (
	"bitlytest/pkg/config"
	"bitlytest/pkg/models"
	"database/sql"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

const (
	tableName = "bitlytest"
)

func (s *Storage) Insert(url models.Url) (uint16, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert(tableName).Columns("small_url", "origin_url").Values(url.SmallUrl, url.OriginUrl).Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var id uint16
	err = s.db.QueryRow(query, args...).Scan(&id)

	return id, err
}

func (s *Storage) Update(url models.Url) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(tableName).Set("small_url", url.SmallUrl).Set("origin_url", url.OriginUrl).Where(squirrel.Eq{"id": url.Id}).ToSql()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = s.db.Exec(query, args...)
	return err
}

func (s *Storage) Delete(url models.Url) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Delete(tableName).Where(squirrel.Eq{"id": url.Id}).ToSql()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = s.db.Exec(query, args...)
	return err
}

func (s *Storage) Get() ([]models.Url, error) {
	query, _, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("id", "small_url", "origin_url").From(tableName).ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	urls := []models.Url{}
	err = s.db.Select(&urls, query)

	/*rows, err := s.db.Query(query)*/

	if err != nil {
		log.Println(err)
		return nil, err
	}

	/*var (
		id         uint16
		small_url  string
		origin_url string
	)*/

	/*defer rows.Close()

	urls := []models.Url{}
	for rows.Next() {
		err := rows.Scan(&id, &small_url, &origin_url)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		urls = append(urls, models.Url{
			Id:        id,
			SmallUrl:  small_url,
			OriginUrl: origin_url,
		})

		log.Println(id, small_url, origin_url)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	*/
	return urls, nil
}

func (s *Storage) GetBySmallUrl(url models.Url) (models.Url, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("id", "small_url", "origin_url").From(tableName).Where(squirrel.Eq{"small_url": url.SmallUrl}).ToSql()
	if err != nil {
		log.Println(err)
		return models.Url{}, err
	}
	/*var (
		id         uint16
		small_url  string
		origin_url string
	)*/

	//err = s.db.QueryRow(query, args...).Scan(&id, &small_url, &origin_url)
	url = models.Url{}
	err = s.db.Get(&url, query, args...)

	if err == sql.ErrNoRows {
		return models.Url{}, models.ErrNotFound
	}

	return url, err
	/*return models.Url{
		Id:        id,
		SmallUrl:  small_url,
		OriginUrl: origin_url,
	}, err*/
}

func DBConnect(cfg config.DbConfig) *sqlx.DB {
	//читать из dbconfig?
	//пока нет
	psql := fmt.Sprintf("postgresql://postgres:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Host, cfg.DbPort, cfg.DbName, cfg.Sslmode)
	db, err := sqlx.Connect("postgres", psql)
	if err != nil {
		log.Println(err)
	}
	return db
}
