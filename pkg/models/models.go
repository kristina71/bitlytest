package models

import "errors"

type Url struct {
	Id        uint16 `json:"id" db:"id"`
	SmallUrl  string `json:"small_url" db:"small_url"`
	OriginUrl string `json:"origin_url" db:"origin_url"`
}

var ErrNotFound = errors.New("not found")
