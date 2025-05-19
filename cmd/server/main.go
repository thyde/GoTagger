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

	// Seed the database if empty
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM bookmarks").Scan(&count)
	if err != nil {
		log.Fatalf("failed to check bookmarks count: %v", err)
	}
	if count == 0 {
		if err := db.SeedDummyData(database); err != nil {
			log.Fatalf("failed to seed dummy data: %v", err)
		}
	}

	// Remove duplicate bookmarks (same url), keeping the one with the lowest id
	_, err = database.Exec(`
	DELETE FROM bookmarks
	WHERE id NOT IN (
		SELECT MIN(id) FROM (
			SELECT id FROM bookmarks GROUP BY url
		)
	);
	`)
	if err != nil {
		log.Fatalf("failed to remove duplicate bookmarks: %v", err)
	}

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, database)
	mux.Handle("/", http.FileServer(http.Dir("./cmd/server/static")))

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", mux)
}
