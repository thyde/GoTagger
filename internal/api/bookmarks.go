package api

import (
	"GoTagger/internal/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

// Handler to list all bookmarks
func ListBookmarksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, url, title, tags, favorite, created_at, updated_at FROM bookmarks")
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bookmarks []model.Bookmark
		for rows.Next() {
			var b model.Bookmark
			var tags string
			if err := rows.Scan(&b.ID, &b.URL, &b.Title, &tags, &b.Favorite, &b.CreatedAt, &b.UpdatedAt); err != nil {
				http.Error(w, "Scan error", http.StatusInternalServerError)
				return
			}
			b.Tags = parseTags(tags)
			bookmarks = append(bookmarks, b)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookmarks)
	}
}

// Handler to add a new bookmark
func AddBookmarkHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b model.Bookmark
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		tags := formatTags(b.Tags)
		_, err := db.Exec("INSERT INTO bookmarks (url, title, tags, favorite) VALUES (?, ?, ?, ?)", b.URL, b.Title, tags, b.Favorite)
		if err != nil {
			http.Error(w, "Insert error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// Handler to search bookmarks by tag or keyword
func SearchBookmarksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		rows, err := db.Query(`SELECT id, url, title, tags, favorite, created_at, updated_at FROM bookmarks WHERE tags LIKE ? OR title LIKE ? OR url LIKE ?`, "%"+q+"%", "%"+q+"%", "%"+q+"%")
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var bookmarks []model.Bookmark
		for rows.Next() {
			var b model.Bookmark
			var tags string
			if err := rows.Scan(&b.ID, &b.URL, &b.Title, &tags, &b.Favorite, &b.CreatedAt, &b.UpdatedAt); err != nil {
				http.Error(w, "Scan error", http.StatusInternalServerError)
				return
			}
			b.Tags = parseTags(tags)
			bookmarks = append(bookmarks, b)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookmarks)
	}
}

// Helper functions for tags
func parseTags(tags string) []string {
	if tags == "" {
		return nil
	}
	return splitAndTrim(tags, ",")
}

func formatTags(tags []string) string {
	return joinAndTrim(tags, ",")
}

func splitAndTrim(s, sep string) []string {
	var out []string
	for _, part := range split(s, sep) {
		trimmed := trim(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func joinAndTrim(parts []string, sep string) string {
	var out []string
	for _, part := range parts {
		trimmed := trim(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return join(out, sep)
}

// Minimal string helpers (replace with strings package in real code)
func split(s, sep string) []string           { return strings.Split(s, sep) }
func trim(s string) string                   { return strings.TrimSpace(s) }
func join(parts []string, sep string) string { return strings.Join(parts, sep) }
