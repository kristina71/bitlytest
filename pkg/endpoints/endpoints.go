package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kristina71/bitlytest/pkg/models"
	"github.com/kristina71/bitlytest/pkg/requestparser"
	"github.com/kristina71/bitlytest/pkg/service"

	"github.com/gorilla/mux"
)

func New(service *service.Service) http.Handler {
	r := mux.NewRouter()

	e := endpoint{service: service}

	r.HandleFunc("/all", e.GetAllUrl).Methods(http.MethodPost)
	r.HandleFunc("/create", e.CreateUrl).Methods(http.MethodPost)
	r.HandleFunc("/delete", e.DeleteUrl).Methods(http.MethodPost)
	r.HandleFunc("/edit", e.UpdateUrl).Methods(http.MethodPost)
	r.HandleFunc("/{small:.*}", e.Get)

	r.Handle("/", http.FileServer(http.Dir("./ui"))).Methods(http.MethodGet)

	staticDir := "/ui/js/"
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	return r
}

type endpoint struct {
	service *service.Service
}

func (e endpoint) Get(w http.ResponseWriter, r *http.Request) {
	url := models.Url{}
	url.SmallUrl = strings.Trim(r.URL.Path, "/")

	url, err := e.service.GetUrl(r.Context(), url)
	if err != nil {
		reportError(err, w)
		return
	}

	http.Redirect(w, r, url.OriginUrl, http.StatusPermanentRedirect)
}

func (e endpoint) GetAllUrl(w http.ResponseWriter, r *http.Request) {
	urls, err := e.service.GetAllUrl(r.Context())
	if err != nil {
		reportError(err, w)
		return
	}

	b, err := json.Marshal(urls)
	if err != nil {
		reportError(err, w)
		return
	}
	w.Write(b)
}

func (e endpoint) CreateUrl(w http.ResponseWriter, r *http.Request) {
	url, resp, err := requestparser.Unmarshal(w, r)
	fmt.Println(string(resp))

	if err != nil {
		reportError(err, w)
		return
	}

	fmt.Println(url)

	url, err = e.service.CreateUrl(r.Context(), url)

	if err != nil {
		reportError(err, w)
		return
	}

	b, err := json.Marshal(url)
	if err != nil {
		reportError(err, w)
		return
	}

	w.Write(b)
}

func (e endpoint) UpdateUrl(w http.ResponseWriter, r *http.Request) {
	url, _, err := requestparser.Unmarshal(w, r)

	if err != nil {
		reportError(err, w)
		return
	}

	fmt.Println(url)

	url, err = e.service.UpdateUrl(r.Context(), url)
	if err != nil {
		reportError(err, w)
		return
	}

	b, err := json.Marshal(url)
	if err != nil {
		reportError(err, w)
		return
	}
	w.Write(b)
}

func (e endpoint) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	url, _, err := requestparser.Unmarshal(w, r)

	if err != nil {
		reportError(err, w)
		return
	}
	fmt.Println(url)
	fmt.Println(url.Id)

	err = e.service.DeleteUrl(r.Context(), url)

	reportError(err, w)
}

func reportError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Println(err)
		log.Printf("%+v\n", err)
		fmt.Printf("%T\n", err)
		switch {
		case errors.As(err, &models.NotFound{}):
			http.Error(w, err.Error(), http.StatusNotFound)
		/*case errors.Is(err, models.ErrCannotUnmarshal):
		http.Error(w, err.Error(), http.StatusNotFound)*/
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}
