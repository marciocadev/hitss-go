package api

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
		if !cpfValidate(c.CPF) {
			http.Error(w, "CPF inválido", http.StatusBadRequest)
			return
		}
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

func removeSpecialCharacters(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

func cpfValidate(cpf string) bool {
	// CPF deve ter 11 dígitos
	if len(cpf) != 11 {
		return false
	}

	// Verifique se todos os dígitos são iguais, o que tornaria o CPF inválido
	if cpf == "00000000000" || cpf == "11111111111" || cpf == "22222222222" ||
		cpf == "33333333333" || cpf == "44444444444" || cpf == "55555555555" ||
		cpf == "66666666666" || cpf == "77777777777" || cpf == "88888888888" ||
		cpf == "99999999999" {
		return false
	}

	// valida primeiro dígito
	cpfd1 := cpf[:9]
	soma := 0
	for i := 0; i < len(cpfd1); i++ {
		num, err := strconv.Atoi(string(cpfd1[i]))
		if err != nil {
			return false
		}
		soma += num * (i + 1)
	}
	d1 := soma % 11
	if d1 == 10 {
		d1 = 0
	}

	// calcula segundo dígito
	cpfd2 := cpf[:10]
	soma = 0
	for i := 0; i < len(cpfd2); i++ {
		num, err := strconv.Atoi(string(cpfd2[i]))
		if err != nil {
			return false
		}
		soma += num * (i)
	}
	d2 := soma % 11
	if d2 == 10 {
		d2 = 0
	}

	// se o primeiro dígito é diferente retorna false
	if string(cpf[9]) != strconv.Itoa(d1) {
		return false
	}
	// se o segundo dígito é diferente retorna false
	if string(cpf[10]) != strconv.Itoa(d2) {
		return false
	}

	// Se passou por todas as verificações, o CPF é válido
	return true
}
