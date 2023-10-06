package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/cliente", s.createClient()).Methods("POST")
}

func removeSpecialCharacters(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

func createID(s string) string {
	hasher := md5.New()
	io.WriteString(hasher, s)
	hashBytes := hasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	return hashString
}

func (s *Server) createClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c Client
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c.CPF = removeSpecialCharacters(c.CPF)
		c.ID = createID(c.CPF)

		StartPublishing(c)

		result := fmt.Sprintf("Cliente %s enviado para o cadastro", c.ID)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
