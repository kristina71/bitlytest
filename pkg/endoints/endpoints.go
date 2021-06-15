package endoints

/*import (
	"bitlytest/pkg/models"
	"bitlytest/pkg/service"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)*/

/*func Delete11Url(w http.ResponseWriter, r *http.Request) {
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

	err := s.Delete1Url(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}*/
