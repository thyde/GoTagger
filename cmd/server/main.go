// gomark/cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"GoTagger/internal/api"
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

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, database)
	mux.Handle("/", http.FileServer(http.Dir("./cmd/server/static")))

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", mux)
}
