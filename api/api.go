package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/rs/cors"
)

var getLink = regexp.MustCompile(`\/api\/shorten\/*$`)
var resolvLink = regexp.MustCompile(`[A-Za-z0-9]+\/*$`)

type apiServer struct {
	listenAddr string
	store      Storage
}

func newapiServer(listenAddr string, store Storage) *apiServer {
	return &apiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (h *apiServer) Run() {
	router := http.NewServeMux()

	router.Handle("/api/shorten/", h)
	router.Handle("/api/shorten", h)
	router.Handle("/s/", h)

	http.ListenAndServe(h.listenAddr, cors.AllowAll().Handler(router))
}

func (h *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	switch {

	case r.Method == http.MethodGet && resolvLink.MatchString(r.URL.Path):
		log.Println("HTTP Method: GET")
		h.resolve(w, r)
		return

	case r.Method == http.MethodPost && getLink.MatchString(r.URL.Path):
		log.Println("HTTP Method: POST")
		h.shorten(w, r)
		return

	default:
		h.notFound(w, r)
		return
	}
}

func (h *apiServer) shorten(w http.ResponseWriter, r *http.Request) {
	newData := new(link)
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newData)
	newData.ID = strconv.FormatInt(time.Now().Unix(), 16)

	if _, err := url.Parse(newData.URL); err != nil {
		log.Println("ERROR:")
		log.Println(err)
		return
	}

	err := h.store.inserttoDB(newData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "please check the logs"})
		log.Println("ERROR:")
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": newData.ID})
}

func (h *apiServer) resolve(w http.ResponseWriter, r *http.Request) {
	resp := new(link)
	id := resolvLink.FindStringSubmatch(r.URL.Path)[0]
	err := h.store.getfromDB(id, resp)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "url not found"})
		return
	}

	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
}

func (h *apiServer) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "page does not exist"})
}
