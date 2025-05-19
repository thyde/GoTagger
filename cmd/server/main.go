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

	// Remove duplicate bookmarks (same url), keeping the one with the lowest id
	_, err = database.Exec(`
	DELETE FROM bookmarks
	WHERE id NOT IN (
		SELECT MIN(id) FROM bookmarks GROUP BY url
	);
	`)
	if err != nil {
		log.Fatalf("failed to remove duplicate bookmarks: %v", err)
	}

	// Ensure required seed bookmarks
	seedData := []struct {
		url   string
		title string
		tags  string
	}{
		{"https://google.com", "Google", "seed,search"},
		{"https://bing.com", "Bing", "seed,search"},
		{"https://yahoo.com", "Yahoo", "seed,search"},
		{"https://amazon.com", "Amazon", "seed,shopping"},
	}
	for _, s := range seedData {
		var exists int
		err := database.QueryRow("SELECT COUNT(*) FROM bookmarks WHERE url = ?", s.url).Scan(&exists)
		if err != nil {
			log.Fatalf("failed to check for seed url %s: %v", s.url, err)
		}
		if exists == 0 {
			_, err := database.Exec("INSERT INTO bookmarks (url, title, tags, favorite) VALUES (?, ?, ?, ?)", s.url, s.title, s.tags, false)
			if err != nil {
				log.Fatalf("failed to insert seed url %s: %v", s.url, err)
			}
		}
	}

	// Ensure at least 5 bookmarks
	var total int
	err = database.QueryRow("SELECT COUNT(*) FROM bookmarks").Scan(&total)
	if err != nil {
		log.Fatalf("failed to count bookmarks: %v", err)
	}
	for i := 1; total < 5; i++ {
		url := fmt.Sprintf("https://example%d.com", i)
		var exists int
		database.QueryRow("SELECT COUNT(*) FROM bookmarks WHERE url = ?", url).Scan(&exists)
		if exists == 0 {
			_, err := database.Exec("INSERT INTO bookmarks (url, title, tags, favorite) VALUES (?, ?, ?, ?)", url, fmt.Sprintf("Example %d", i), "seed,example", false)
			if err != nil {
				log.Fatalf("failed to insert extra seed: %v", err)
			}
			total++
		}
	}

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, database)
	mux.Handle("/", http.FileServer(http.Dir("./cmd/server/static")))

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", mux)
}
