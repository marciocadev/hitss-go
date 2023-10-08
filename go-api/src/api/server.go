package api

import (
	"crypto/md5"
	"database/sql"
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

func NewServer(db *sql.DB) *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}
	s.routes(db)
	return s
}

func (s *Server) routes(db *sql.DB) {
	s.HandleFunc("/cliente", s.createClient()).Methods("POST")
	s.HandleFunc("/cliente", s.getAllClients(db)).Methods("GET")
	s.HandleFunc("/cliente/{id}", s.getClientByID(db)).Methods("GET")
	s.HandleFunc("/cliente/cpf/{cpf}", s.getClientByCPF(db)).Methods("GET")
	s.HandleFunc("/cliente/{id}", s.removeClient()).Methods("DELETE")
	s.HandleFunc("/cliente/{id}", s.updateClient()).Methods("PUT")
}

func (s *Server) getClientByCPF(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cpfStr, _ := mux.Vars(r)["cpf"]
		cpfStr = removeSpecialCharacters(cpfStr)
		client := GetClientByCPF(db, cpfStr)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(client); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getClientByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["id"]
		client := GetClientByID(db, id)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(client); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) getAllClients(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		todos := GetAllClients(db)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
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

		PublishingInsertNewClient(c)

		result := fmt.Sprintf("Cliente %s enviado para o cadastro", c.ID)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) removeClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["id"]

		PublishingDeleteClient(id)

		result := fmt.Sprintf("Cliente %s enviado para ser removido", id)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) updateClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["id"]

		var c Client
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		PublishingUpdateClient(id, c)

		result := fmt.Sprintf("Cliente %s enviado para ser atualizado", id)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
