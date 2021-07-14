package repositories

import (
	"context"

	"github.com/kristina71/bitlytest/pkg/adapters"
	"github.com/kristina71/bitlytest/pkg/generator"
	"github.com/kristina71/bitlytest/pkg/models"
	"github.com/kristina71/bitlytest/pkg/urlvalidator"
)

type Urls struct {
	adapter *adapters.Storage
}

func New(adapter *adapters.Storage) *Urls {
	return &Urls{adapter: adapter}
}

func (u *Urls) Insert(ctx context.Context, url models.Url) (uint16, error) {
	return u.adapter.Insert(ctx, url)
}

func (u *Urls) GetBySmallUrl(ctx context.Context, url models.Url) (models.Url, error) {
	return u.adapter.GetBySmallUrl(ctx, url)
}

func (u *Urls) Get(ctx context.Context) ([]models.Url, error) {
	return u.adapter.Get(ctx)
}

func (u *Urls) Update(ctx context.Context, url models.Url) error {
	return u.adapter.Update(ctx, url)
}

func (u *Urls) Delete(ctx context.Context, url models.Url) error {
	return u.adapter.Delete(ctx, url)
}

func (u *Urls) GenerateUrl(_ context.Context) string {
	return generator.RandomString()
}

func (u *Urls) ValidateUrl(ctx context.Context, url string) bool {
	return urlvalidator.ValidateUrl(ctx, url)
}
