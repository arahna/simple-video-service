package handlers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func TestList(t *testing.T) {
	repo, cleanup := newVideoRepository()
	defer cleanup()
	r, err := http.NewRequest("GET", "/api/v1/list", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	Router(repo).ServeHTTP(w, r)
 	response := w.Result()

 	if response.StatusCode != http.StatusOK {
 		t.Fatalf("Status code is wrong. Have: %d, want: %d.", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	items := make([]videoListItem, 10)
	if err = json.Unmarshal(jsonString, &items); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
}
