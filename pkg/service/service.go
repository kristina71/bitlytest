package service

import (
	"bitlytest/pkg/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
)

type Storage interface {
	Insert(url models.Url) (uint16, error)
	Update(url models.Url) error
	Delete(url models.Url) error
	Get() ([]models.Url, error)
	GetBySmallUrl(url models.Url) (models.Url, error)
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s Service) CreateUrl(w http.ResponseWriter, r *http.Request) {
	url := models.Url{}
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = json.Unmarshal(resp, &url)
	if err != nil {
		http.Error(w, "Cannot unmarshal json", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	url.SmallUrl = strings.Trim(url.SmallUrl, " ")
	url.SmallUrl = strings.Trim(url.SmallUrl, "/")
	if url.SmallUrl == "" {
		http.Error(w, "empty url", http.StatusBadRequest)
		log.Println("empty url")
		return
	}

	url.OriginUrl = strings.Trim(url.OriginUrl, " ")
	url.OriginUrl = strings.Trim(url.OriginUrl, "/")
	if url.OriginUrl == "" {
		http.Error(w, "empty origin url", http.StatusBadRequest)
		log.Println("empty origin url")
		return
	}

	if !IsUrl(url.OriginUrl) {
		http.Error(w, "incorrect origin url", http.StatusBadRequest)
		log.Println("incorrect origin url")
		return
	}

	fmt.Println(url)

	url.Id, err = s.storage.Insert(url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	b, err := json.Marshal(url)
	if err != nil {
		http.Error(w, "Cannot marshal json", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Write(b)
}

func (s Service) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	url := models.Url{}
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = json.Unmarshal(resp, &url)
	if err != nil {
		http.Error(w, "Cannot unmarshal json", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	fmt.Println(url)
	fmt.Println(url.Id)

	err = s.storage.Delete(url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (s Service) Delete1Url(url models.Url) error {
	return s.storage.Delete(url)
}

func (s Service) UpdateUrl(w http.ResponseWriter, r *http.Request) {
	url := models.Url{}
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = json.Unmarshal(resp, &url)
	if err != nil {
		http.Error(w, "Cannot unmarshal json", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	url.SmallUrl = strings.Trim(url.SmallUrl, " ")
	url.SmallUrl = strings.Trim(url.SmallUrl, "/")
	if url.SmallUrl == "" {
		http.Error(w, "empty url", http.StatusBadRequest)
		log.Println("empty url")
		return
	}

	url.OriginUrl = strings.Trim(url.OriginUrl, " ")
	url.OriginUrl = strings.Trim(url.OriginUrl, "/")
	if url.SmallUrl == "" {
		http.Error(w, "empty origin url", http.StatusBadRequest)
		log.Println("empty origin url")
		return
	}
	if !IsUrl(url.OriginUrl) {
		http.Error(w, "incorrect origin url", http.StatusInternalServerError)
		log.Println("incorrect origin url")
		return
	}

	fmt.Println(url)

	err = s.storage.Update(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Write(resp)
}

func (s Service) GetUrl(w http.ResponseWriter, r *http.Request) {
	url := models.Url{}
	url.SmallUrl = strings.Trim(r.URL.Path, "/")
	//url.SmallUrl = mux.Vars(r)["small_url"]

	fmt.Println(url)

	urls, err := s.storage.GetBySmallUrl(url)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		log.Println(err)
		return
	}

	http.Redirect(w, r, urls.OriginUrl, http.StatusMovedPermanently)
}

func (s Service) GetAllUrl(w http.ResponseWriter, r *http.Request) {
	urls, err := s.storage.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	b, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, "Cannot marshal json", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Write(b)
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
