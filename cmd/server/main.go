// gomark/cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"GoTagger/internal/db"
)

func main() {
	database, err := db.NewSQLiteDB("")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	err = db.Migrate(database)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	err = db.SeedDummyData(database)
	if err != nil {
		log.Fatalf("failed to seed dummy data: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to GoMark!")
	})

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
