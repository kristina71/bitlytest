package adapters

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/kristina71/bitlytest/pkg/config"
	"github.com/kristina71/bitlytest/pkg/models"
	"github.com/pkg/errors"

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

func (s *Storage) Insert(ctx context.Context, url models.Url) (uint16, error) {
	if url.CreatedAt.IsZero() {
		url.CreatedAt = time.Now().UTC()
	}
	if url.UpdateAt.IsZero() {
		url.UpdateAt = time.Now().UTC()
	}

	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert(tableName).Columns("small_url", "origin_url", "created_at", "updated_at").Values(url.SmallUrl, url.OriginUrl, url.CreatedAt, url.UpdateAt).Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var id uint16
	err = s.db.QueryRow(query, args...).Scan(&id)

	return id, err
}

func (s *Storage) Update(ctx context.Context, url models.Url) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(tableName).Set("small_url", url.SmallUrl).Set("origin_url", url.OriginUrl).Where(squirrel.Eq{"id": url.Id}).ToSql()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = s.db.Exec(query, args...)
	return err
}

func (s *Storage) Delete(ctx context.Context, url models.Url) error {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Delete(tableName).Where(squirrel.Eq{"id": url.Id}).ToSql()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = s.db.Exec(query, args...)
	return err
}

func (s *Storage) Get(ctx context.Context) ([]models.Url, error) {
	query, _, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("id", "small_url", "origin_url").From(tableName).ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	urls := []models.Url{}
	err = s.db.Select(&urls, query)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return urls, nil
}

func (s *Storage) GetBySmallUrl(ctx context.Context, url models.Url) (models.Url, error) {
	query, args, err := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("id", "small_url", "origin_url").From(tableName).Where(squirrel.Eq{"small_url": url.SmallUrl}).ToSql()
	if err != nil {
		log.Println(err)
		return models.Url{}, err
	}

	url = models.Url{}
	err = s.db.Get(&url, query, args...)

	if err == sql.ErrNoRows {
		return models.Url{}, errors.WithStack(models.NotFoundError())
	}

	return url, err
}

func DBConnect(cfg config.Cfg) *sqlx.DB {
	db, err := sqlx.Connect(cfg.DbDialect, cfg.DbDsn)
	if err != nil {
		log.Println(err)
	}
	return db
}
