package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func OpenConn() *sql.DB {
	// connection string
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"postgres", 5432, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	// open database
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to open database", err)
	}
	return db
}

func convertDateFormat(s string) string {
	layout := "2006-01-02T15:04:05Z"
	dt, err := time.Parse(layout, s)
	if err != nil {
		log.Fatalf("%s: %v", "Date string format fail", err)
	}
	outputLayout := "02/01/2006"
	return dt.Format(outputLayout)
}

func GetAllClients(db *sql.DB) []ClientDTO {
	// create statement
	query := "SELECT id, nome, sobrenome, contato, endereco, dtnascimento FROM hitss.cliente"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to create the statement", err)
	}
	defer rows.Close()

	var todos []ClientDTO
	for rows.Next() {
		var c ClientDTO
		err := rows.Scan(&c.ID, &c.Nome, &c.Sobrenome, &c.Contato, &c.Endereco, &c.DtNascimento)
		if err != nil {
			log.Fatalf("%s: %v", "Data can not be retrieved", err)
		}
		c.DtNascimento = convertDateFormat(c.DtNascimento)
		todos = append(todos, c)
	}

	return todos
}

func GetClientByID(db *sql.DB, id string) ClientDTO {
	// create statement
	query := "SELECT id, nome, sobrenome, contato, endereco, dtnascimento FROM hitss.cliente WHERE id = $1"
	rows := db.QueryRow(query, id)
	var c ClientDTO
	err := rows.Scan(&c.ID, &c.Nome, &c.Sobrenome, &c.Contato, &c.Endereco, &c.DtNascimento)
	if err != nil {
		log.Fatalf("%s: %v", "Data can not be retrieved", err)
	}
	c.DtNascimento = convertDateFormat(c.DtNascimento)
	return c
}

func GetClientByCPF(db *sql.DB, cpf string) ClientDTO {
	// create statement
	query := "SELECT id, nome, sobrenome, contato, endereco, dtnascimento FROM hitss.cliente WHERE cpf = $1"
	rows := db.QueryRow(query, cpf)
	var c ClientDTO
	err := rows.Scan(&c.ID, &c.Nome, &c.Sobrenome, &c.Contato, &c.Endereco, &c.DtNascimento)
	if err != nil {
		log.Fatalf("%s: %v", "Data can not be retrieved", err)
	}
	c.DtNascimento = convertDateFormat(c.DtNascimento)
	return c
}
