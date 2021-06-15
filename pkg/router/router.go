package router

import (
	"bitlytest/pkg/service"
	"bitlytest/pkg/storage"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) http.Handler {
	storage := storage.New(db)
	s := service.New(storage)
	router := mux.NewRouter()

	router.HandleFunc("/all", s.GetAllUrl).Methods(http.MethodPost)
	router.HandleFunc("/create", s.CreateUrl).Methods(http.MethodPost)
	router.HandleFunc("/delete", s.DeleteUrl).Methods(http.MethodPost)
	router.HandleFunc("/edit", s.UpdateUrl).Methods(http.MethodPost)
	router.HandleFunc("/{small_url:.+}", s.GetUrl).Methods(http.MethodGet)

	router.Handle("/", http.FileServer(http.Dir("./ui"))).Methods(http.MethodGet)

	staticDir := "/ui/js/"
	router.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	return router
}
