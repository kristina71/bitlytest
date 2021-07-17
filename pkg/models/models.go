package models

import (
	"time"
)

type Url struct {
	Id        uint16    `json:"id" db:"id"`
	SmallUrl  string    `json:"small_url" db:"small_url"`
	OriginUrl string    `json:"origin_url" db:"origin_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdateAt  time.Time `json:"updated_at" db:"updated_at"`
}
