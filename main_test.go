package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lithammer/shortuuid"
	"github.com/marcogregorius/url-shortener/app"
)

var a app.App

func TestMain(m *testing.M) {
	a.Initialize()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	a.DB.Exec("DELETE FROM tb_shortlinks")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		if len(message) == 0 {
			message = fmt.Sprintf("%v != %v", a, b)
		}
		t.Fatal(message)
	}
}

func assertNotNil(t *testing.T, a interface{}, message string) {
	if a == nil {
		if len(message) == 0 {
			message = fmt.Sprintf("nil unexpected")
		}
		t.Fatal(message)
	}
}

func TestCreateShortlink(t *testing.T) {
	clearTable()

	url := `https://blog.golang.org/`
	var jsonStr = bytes.Replace([]byte(`{"source_url": "URL"}`), []byte(`URL`), []byte(url), 1)
	req, _ := http.NewRequest("POST", "/api/shortlinks", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	notNillFields := []string{"id", "created_at"}
	for _, v := range notNillFields {
		assertNotNil(t, m[v], "")
	}

	assertEqual(t, m["source_url"], string(url), "")
}
func TestCreateShortlinkInvalidURL(t *testing.T) {
	clearTable()

	// try to create shortlink from invalid URL format
	url := `"something_randomDotCom"`
	var jsonStr = bytes.Replace([]byte(`{"source_url": URL}`), []byte("URL"), []byte(url), 1)
	req, _ := http.NewRequest("POST", "/api/shortlinks", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestGetShortlink(t *testing.T) {
	clearTable()

	id, sourceUrl := addShortlink()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/shortlinks/%v", id), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	assertEqual(t, m["id"], id, "")
	assertEqual(t, m["source_url"], sourceUrl, "")
}

func TestGetShortlinkNotFound(t *testing.T) {
	clearTable()

	id := "some_random_id"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/shortlinks/%v", id), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addShortlink() (id string, sourceUrl string) {
	id = shortuuid.New()[:8]
	sourceUrl = "https://www.google.com"
	a.DB.Exec("INSERT INTO tb_shortlinks(id,source_url,created_at) VALUES ($1, $2, $3)",
		id, sourceUrl, time.Now())
	return
}

func TestVisitShortlink(t *testing.T) {
	clearTable()

	id, sourceUrl := addShortlink()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", id), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusMovedPermanently, response.Code)
	assertEqual(
		t,
		strings.Contains(response.Body.String(), sourceUrl),
		true,
		fmt.Sprintf("Redirection URL not found. Expected to contain %v", sourceUrl),
	)

	// check if visited has increased by 1
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/shortlinks/%v", id), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	assertEqual(t, int(m["visited"].(float64)), 1, "")
}

func TestVisitShortlinkNotFound(t *testing.T) {
	clearTable()

	id := "some_random_id"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", id), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
