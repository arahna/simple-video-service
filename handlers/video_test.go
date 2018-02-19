package handlers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/arahna/simple-video-service/database"
	"log"
	"database/sql"
)

func TestVideo(t *testing.T) {
	db := initDB()
	defer db.Close()
	r, err := getRequest("sldjfl34-dfgj-523k-jk34-5jk3j45klj34")
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

	item := videoItem{}
	if err = json.Unmarshal(jsonString, &item); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
}

func TestVideoNotFound(t *testing.T) {
	db := initDB()
	defer db.Close()
	r, err := getRequest("non-existent-video")
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	Router(db).ServeHTTP(w, r)
	response := w.Result()

	testStatusCode(response.StatusCode, http.StatusNotFound, t)
}

func testStatusCode(have int, want int, t *testing.T) {
	if have != want {
		t.Fatalf("Status code is wrong. Have: %d, want: %d.", have, want)
	}
}

func getRequest(id string) (*http.Request, error) {
	url := fmt.Sprintf("/api/v1/video/%s", id)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func initDB() *sql.DB {
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
