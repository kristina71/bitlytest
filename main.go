package main

import (
	"bitlytest/pkg/config"
	"bitlytest/pkg/router"
	"bitlytest/pkg/storage"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()
	db := storage.DBConnect(cfg)

	srv := &http.Server{
		Handler:      router.NewRouter(db),
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	defer db.Close()
}
