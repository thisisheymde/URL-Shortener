package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/thisisheymde/URL-shortener/backend/api"
	"github.com/thisisheymde/URL-shortener/backend/storage"
	"github.com/thisisheymde/URL-shortener/backend/types"
)

var store, _ = storage.StartRedis(os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))
var server = api.NewServer(":8081", store)

func TestServer_Shorten(t *testing.T) {
	payload := map[string]string{"url": "http://example.com"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))

	respr := httptest.NewRecorder() // recording response
	server.ServeHTTP(respr, req)

	if respr.Code == 400 {
		t.Errorf("expected 200, got %v", respr.Code)
	}

	body, _ = io.ReadAll(respr.Body)

	newData := new(types.Link)
	json.Unmarshal(body, &newData)

	if len(newData.ID) == 0 {
		t.Errorf("expected shortened code, got %v", string(body))
	}
}

func TestServer_Resolve(t *testing.T) {
	payload := map[string]string{"url": "http://example.com"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))

	respr := httptest.NewRecorder() // recording response
	server.ServeHTTP(respr, req)

	body, _ = io.ReadAll(respr.Body)

	newData := new(types.Link)
	json.Unmarshal(body, &newData) // get the code

	req, _ = http.NewRequest(http.MethodGet, "/s/"+newData.ID, nil)

	respr = httptest.NewRecorder() // recording response
	server.ServeHTTP(respr, req)

	if respr.Code != 303 {
		t.Errorf("expected 303, got %v", respr.Code)
	}

	if respr.Header().Get("Location") != "http://example.com" {
		t.Errorf("expected location http://example.com, got %v", respr.Header().Get("Location"))
	}
}

func TestServer_Put(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPut, "", nil)
	respr := httptest.NewRecorder() // recording response
	server.ServeHTTP(respr, req)

	if respr.Code != 405 {
		t.Errorf("expected 405, got %v", respr.Code)
	}
}

func TestServer_Patch(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPatch, "", nil)
	respr := httptest.NewRecorder() // recording response
	server.ServeHTTP(respr, req)

	if respr.Code != 405 {
		t.Errorf("expected 405, got %v", respr.Code)
	}
}

func TestServer_Delete(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "", nil)
	respr := httptest.NewRecorder() // recording response
	server.ServeHTTP(respr, req)

	if respr.Code != 405 {
		t.Errorf("expected 405, got %v", respr.Code)
	}
}

func TestServer_NotFound(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/randomfoundnot", nil)
	respr := httptest.NewRecorder() // recording response
	server.ServeHTTP(respr, req)

	if respr.Code != 404 {
		t.Errorf("expected 404, got %v", respr.Code)
	}
}

func TestServer_RateLimiting(t *testing.T) {
	requests := 10 // 13*2 = 26, 26 requests
	delay := time.Second

	for i := 0; i < requests; i++ {
		payload := map[string]string{"url": "http://example" + strconv.Itoa(i) + ".com"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))

		respr := httptest.NewRecorder() // recording response
		server.ServeHTTP(respr, req)

		body, _ = io.ReadAll(respr.Body)

		newData := new(types.Link)
		json.Unmarshal(body, &newData) // get the code

		req, _ = http.NewRequest(http.MethodGet, "/s/"+newData.ID, nil)

		respr = httptest.NewRecorder() // recording response
		server.ServeHTTP(respr, req)

		if i == 13 && respr.Code != 429 {
			t.Errorf("Expected 429 on 26th request, too many requests, got %v", respr.Code)
		}

		time.Sleep(delay)
	}
}
