package generator

import (
	"math/rand"
	"sync"
	"time"
)

var once sync.Once

func RandomString() string {
	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})

	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	newUrl := make([]rune, 16)
	for i := range newUrl {
		newUrl[i] = letters[rand.Intn(len(letters))]
	}
	return string(newUrl)
}
