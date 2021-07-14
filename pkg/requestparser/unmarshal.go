package requestparser

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kristina71/bitlytest/pkg/models"
)

func Unmarshal(w http.ResponseWriter, r *http.Request) (models.Url, []byte, error) {
	url := models.Url{}
	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return url, nil, err
	}

	err = json.Unmarshal(resp, &url)
	if err != nil {
		return url, nil, err
	}

	return url, resp, nil
}
