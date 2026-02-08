package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"api/handlers"

	_ "github.com/lib/pq"
)

func ExecuteDDL(db *sql.DB, path string) error {
	ddlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read DDL file: %w", err)
	}

	ddl := string(ddlBytes)

	_, err = db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("cannot execute DDL: %w", err)
	}

	return nil
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@db:5432/myapp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	API := handlers.New(db)

	http.HandleFunc("/api/health", healthHandler)

	http.Handle("/api/register", handlers.EnableCors(http.HandlerFunc(API.RegisterHandler)))
	http.Handle("/api/login", handlers.EnableCors(http.HandlerFunc(API.LoginHandler)))

	http.Handle("/api/load_all_messages", handlers.EnableCors(http.HandlerFunc(API.LoadAllMessagesHandler)))
	http.Handle("/ws", handlers.EnableCors(http.HandlerFunc(API.WsHandler)))

	log.Println("API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
