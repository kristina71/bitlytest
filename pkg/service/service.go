package service

import (
	"context"
	"errors"
	"strings"

	"github.com/kristina71/bitlytest/pkg/models"

	_ "github.com/lib/pq"
)

type Repository interface {
	Insert(ctx context.Context, url models.Url) (uint16, error)
	Update(ctx context.Context, url models.Url) error
	Delete(ctx context.Context, url models.Url) error
	Get(ctx context.Context) ([]models.Url, error)
	GetBySmallUrl(ctx context.Context, url models.Url) (models.Url, error)
	GenerateUrl(ctx context.Context) string
	ValidateUrl(ctx context.Context, url string) bool
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) CreateUrl(ctx context.Context, url models.Url) (models.Url, error) {
	url = trimUrl(url)

	if !s.repo.ValidateUrl(ctx, url.OriginUrl) {
		return url, errors.New("invalid origin url")
	}

	if url.SmallUrl == "" {
		url.SmallUrl = s.repo.GenerateUrl(ctx)
	}

	var err error
	url.Id, err = s.repo.Insert(ctx, url)
	return url, err
}

func (s Service) DeleteUrl(ctx context.Context, url models.Url) error {
	return s.repo.Delete(ctx, url)
}

func (s Service) UpdateUrl(ctx context.Context, url models.Url) (models.Url, error) {
	url = trimUrl(url)

	if !s.repo.ValidateUrl(ctx, url.OriginUrl) {
		return url, errors.New("invalid origin url")
	}

	if url.SmallUrl == "" {
		url.SmallUrl = s.repo.GenerateUrl(ctx)
	}

	err := s.repo.Update(ctx, url)
	return url, err

}

func (s Service) GetUrl(ctx context.Context, url models.Url) (models.Url, error) {
	return s.repo.GetBySmallUrl(ctx, url)
}

func (s Service) GetAllUrl(ctx context.Context) ([]models.Url, error) {
	return s.repo.Get(ctx)
}

func trimUrl(url models.Url) models.Url {
	url.SmallUrl = strings.Trim(url.SmallUrl, " ")
	url.SmallUrl = strings.Trim(url.SmallUrl, "/")

	url.OriginUrl = strings.Trim(url.OriginUrl, " ")
	return url
}
