package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kristina71/bitlytest/pkg/adapters"
	"github.com/kristina71/bitlytest/pkg/config"
	endpoints "github.com/kristina71/bitlytest/pkg/endpoints"
	"github.com/kristina71/bitlytest/pkg/repositories"
	"github.com/kristina71/bitlytest/pkg/service"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()
	db := adapters.DBConnect(cfg)
	adapters := adapters.New(db)
	repo := repositories.New(adapters)
	service := service.New(repo)

	srv := &http.Server{
		Handler:      endpoints.New(service),
		Addr:         fmt.Sprintf("%s:%s", "localhost", "8000"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	defer db.Close()
}
