package handlers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/arahna/simple-video-service/internal/pkg/model"
)

func TestStatus(t *testing.T) {
	db := initDB()
	defer db.Close()
	r, err := getStatusRequest("sldjfl34-dfgj-523k-jk34-5jk3j45klj34")
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	Router(db).ServeHTTP(w, r)
	response := w.Result()

	testStatusCode(response.StatusCode, http.StatusOK, t)

	jsonString, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	item := statusResponse{}
	if err = json.Unmarshal(jsonString, &item); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}

	if item.Status != model.VideoReady {
		t.Errorf("Video status is wrong. Have: %d, want: %d.", item.Status, model.VideoReady)
	}
}

func TestStatusVideoNotFound(t *testing.T) {
	db := initDB()
	defer db.Close()
	r, err := getStatusRequest("non-existent-video")
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	Router(db).ServeHTTP(w, r)
	response := w.Result()

	testStatusCode(response.StatusCode, http.StatusNotFound, t)
}

func getStatusRequest(id string) (*http.Request, error) {
	url := fmt.Sprintf("/api/v1/video/%s/status", id)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}