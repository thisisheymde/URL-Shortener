package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/thisisheymde/URL-shortener/backend/types"
	"github.com/thisisheymde/URL-shortener/backend/utils"
)

func (s *Server) shorten(w http.ResponseWriter, r *http.Request) {
	newData := new(types.Link)
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newData)

	if len(newData.URL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "url is required"})
		return
	}

	if len(newData.ID) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	newData.ID = utils.Hashing(newData.URL)[:8]

	// checks if URL is valid or not
	_, err := url.ParseRequestURI(newData.URL)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid URL"})
		return
	}

	err = s.store.InserttoDB(newData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "please check the logs"})
		log.Print("ERROR:")
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": newData.ID})
}

func (s *Server) resolve(w http.ResponseWriter, r *http.Request) {
	resp := new(types.Link)
	var err error
	id := resolvLink.FindStringSubmatch(r.URL.Path)[0]
	id = strings.ReplaceAll(id, "/", "")

	resp.URL, err = s.store.GetfromDB(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "link not found"})
		return
	}

	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "page does not exist"})
}
