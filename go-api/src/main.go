package main

import (
	"hitss/api"
	"net/http"
)

func main() {
	db := api.OpenConn()

	srv := api.NewServer(db)

	certPath := "/work/cert/server-cert.pem"
	keyPath := "/work/cert/server-key.pem"

	http.ListenAndServeTLS(":8080", certPath, keyPath, srv)

	// close database
	defer db.Close()
}
