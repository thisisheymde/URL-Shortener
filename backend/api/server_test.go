package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thisisheymde/URL-shortener/backend/api"
	"github.com/thisisheymde/URL-shortener/backend/storage"
	"github.com/thisisheymde/URL-shortener/backend/types"
)

var store, _ = storage.StartRedis("containers-us-west-33.railway.app:7772", "pOK9WhpY0SZ5TLh8T6ui")
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
	json.Unmarshal(body, &newData)

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

// func TestServer_RateLimiting(t *testing.T) {
// }
