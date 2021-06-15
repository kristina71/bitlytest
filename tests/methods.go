package tests

import (
	"bitlytest/pkg/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func DeleteItem(ts *httptest.Server, url models.Url) error {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(url)
	_, err := ts.Client().Post(ts.URL+"/delete", "application/json", payloadBuf)
	return err
}

func GetItemBySMallUrl(ts *httptest.Server, url models.Url) (*http.Response, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(url)
	res, err := ts.Client().Post(ts.URL+"/delete", "application/json", payloadBuf)
	return res, err
}

func CreateItem(ts *httptest.Server, url models.Url) (*http.Response, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(url)
	res, err := ts.Client().Post(ts.URL+"/create", "application/json", payloadBuf)
	return res, err
}

func EditItem(ts *httptest.Server, url models.Url) (*http.Response, error) {
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(url)
	res, err := ts.Client().Post(ts.URL+"/edit", "application/json", payloadBuf)
	return res, err
}
