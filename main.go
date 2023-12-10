// main.go
package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	// Initialize the database connection during program startup
	var err error
	db, err = ConnectDB()
	if err != nil {
		panic(err)
	}
	initDB()
}

func main() {
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/neighborhoods", GetNeighborhoods).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
