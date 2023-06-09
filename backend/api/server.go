package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/rs/cors"
	"github.com/thisisheymde/URL-shortener/backend/storage"
)

var getLink = regexp.MustCompile(`\/api\/shorten\/*$`)
var resolvLink = regexp.MustCompile(`[A-Za-z0-9]+\/*$`)

type Server struct {
	listenAddr string
	store      storage.Storage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Run() {
	router := http.NewServeMux()
	router.Handle("/api/shorten", s)
	router.Handle("/api/shorten/", s)
	router.Handle("/s/", s)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:8080", "http://localhost:8080"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
	})

	http.ListenAndServe(s.listenAddr, c.Handler(router))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//Rate Limiting
	err := RateLimiting(w, r)

	if errors.Is(err, ErrRateExceeded) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"error": "too many requests"})
		return
	}

	if err != nil {
		log.Print("ERR: ")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "check logs, bad cache"})
		return
	}

	switch {
	case r.Method == http.MethodGet && resolvLink.MatchString(r.URL.Path):
		s.resolve(w, r)
		return

	case r.Method == http.MethodPost && getLink.MatchString(r.URL.Path):
		s.shorten(w, r)
		return

	case r.Method == http.MethodPut || r.Method == http.MethodPatch:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return

	case r.Method == http.MethodDelete:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return

	default:
		s.notFound(w, r)
		return
	}
}
