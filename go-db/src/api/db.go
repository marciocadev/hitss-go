package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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

func GetUpdateStatement(db *sql.DB, c Client) (*sql.Stmt, []interface{}) {
	// update statement
	var sqlQuery string = "UPDATE hitss.cliente SET "
	var i int = 1

	var params []interface{}

	for fieldName, fieldValue := range clientToMap(c) {
		if fieldValue.(string) != "" && fieldName != "ID" {
			sqlQuery += fmt.Sprintf(" %s = $%d, ", fieldName, i)
			i++
			params = append(params, fieldValue)
			log.Printf("%s", fieldValue)
		}
	}

	sqlQuery = strings.TrimSuffix(sqlQuery, ", ")
	sqlQuery += fmt.Sprintf(" WHERE id = $%d", i)
	params = append(params, c.ID)

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to create the statement", err)
	}
	return stmt, params
}

func UpdateClient(stmt *sql.Stmt, params []interface{}) {
	// update client to database
	log.Println("inserting ...")
	_, err := stmt.Exec(params...)

	if err != nil {
		if isDuplicateError(err) {
			// If client already exists, log a message and remove it from RabbitMQ
			log.Print("Cliente já foi cadastrado")
		} else {
			// If fail from any oher way, the message stay in queue to reprocessing
			log.Fatalf("%s: %v", "Insert failed", err)
		}
	}
}

func GetDeleteStatement(db *sql.DB) *sql.Stmt {
	// delete statement
	var sqlQuery string = "DELETE FROM hitss.cliente " +
		"WHERE id = $1"
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to create the statement", err)
	}
	return stmt
}

func DeleteClient(stmt *sql.Stmt, id string) {
	// insert client to database
	log.Println("removing ...")
	res, err := stmt.Exec(id)
	if err != nil || res == nil {
		// If fail from any oher way, the message stay in queue to reprocessing
		log.Fatalf("%s: %v", "Insert failed", err)
	}
}

func GetInsertStatement(db *sql.DB) *sql.Stmt {
	// create statement
	var sqlQuery string = "INSERT INTO hitss.cliente " +
		"(id, nome, sobrenome, contato, endereco, dtNascimento, cpf)" +
		" values ($1, $2, $3, $4, $5, $6, $7)"
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Fatalf("%s: %v", "Failed to create the statement", err)
	}
	return stmt
}

func InsertClient(stmt *sql.Stmt, id string, nome string, sobrenome string, contato string, end string, dt time.Time, cpf string) {
	// insert client to database
	log.Println("inserting ...")
	_, err := stmt.Exec(id, nome, sobrenome, contato, end, dt, cpf)

	if err != nil {
		if isDuplicateError(err) {
			// If client already exists, log a message and remove it from RabbitMQ
			log.Print("Cliente já foi cadastrado")
		} else {
			// If fail from any oher way, the message stay in queue to reprocessing
			log.Fatalf("%s: %v", "Insert failed", err)
		}
	}
}

// isDuplicateError checks if the error is a PostgreSQL unique constraint violation error
func isDuplicateError(err error) bool {
	return err != nil && (containsIgnoreCase(err.Error(), "unique constraint") || containsIgnoreCase(err.Error(), "duplicate key"))
}

// containsIgnoreCase checks if a string contains another string (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// Utility function to convert Client struct to a map
func clientToMap(client Client) map[string]interface{} {
	clientMap := make(map[string]interface{})
	clientMap["ID"] = client.ID
	clientMap["Nome"] = client.Nome
	clientMap["Sobrenome"] = client.Sobrenome
	clientMap["Contato"] = client.Contato
	clientMap["Endereco"] = client.Endereco
	clientMap["DtNascimento"] = client.DtNascimento
	clientMap["CPF"] = client.CPF
	return clientMap
}
