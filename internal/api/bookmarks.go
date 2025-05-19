package api

import (
	"GoTagger/internal/db"
	"GoTagger/internal/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Handler to list bookmarks with pagination and sorting
func ListBookmarksHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("search")
		page := 1
		perPage := 10
		sort := r.URL.Query().Get("sort")
		if v := r.URL.Query().Get("page"); v != "" {
			fmt.Sscanf(v, "%d", &page)
			if page < 1 {
				page = 1
			}
		}
		if v := r.URL.Query().Get("per_page"); v != "" {
			fmt.Sscanf(v, "%d", &perPage)
			if perPage < 1 {
				perPage = 10
			}
		}
		bookmarks, total, err := db.ListBookmarksPaginated(database, q, sort, page, perPage)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"bookmarks": bookmarks,
			"total":     total,
			"page":      page,
			"per_page":  perPage,
		})
	}
}

// Handler to add a new bookmark
func AddBookmarkHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b model.Bookmark
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := db.CreateBookmark(database, &b); err != nil {
			http.Error(w, "Insert error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// Handler to search bookmarks by tag or keyword
func SearchBookmarksHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		bookmarks, err := db.SearchBookmarksByKeywordOrTag(database, q)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bookmarks)
	}
}

// Handler to update an existing bookmark
func UpdateBookmarkHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b model.Bookmark
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if b.ID == 0 {
			http.Error(w, "Missing bookmark ID", http.StatusBadRequest)
			return
		}
		if err := db.UpdateBookmark(database, &b); err != nil {
			http.Error(w, "Update error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Handler to delete a bookmark
func DeleteBookmarkHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Missing bookmark ID", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid bookmark ID", http.StatusBadRequest)
			return
		}
		if err := db.DeleteBookmark(database, id); err != nil {
			http.Error(w, "Delete error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Handler to list all tags and their usage counts
func ListTagsHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.Query(`SELECT tags FROM bookmarks`)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		tagCounts := make(map[string]int)
		for rows.Next() {
			var tags string
			if err := rows.Scan(&tags); err != nil {
				http.Error(w, "Scan error", http.StatusInternalServerError)
				return
			}
			for _, tag := range parseTags(tags) {
				tagCounts[tag]++
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tagCounts)
	}
}

// Handler for exporting bookmarks as JSON
func ExportBookmarksHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookmarks, err := db.ListBookmarks(database)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=bookmarks.json")
		json.NewEncoder(w).Encode(bookmarks)
	}
}

// Handler for importing bookmarks from JSON
func ImportBookmarksHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bookmarks []model.Bookmark
		if err := json.NewDecoder(r.Body).Decode(&bookmarks); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		for _, b := range bookmarks {
			_ = db.CreateBookmark(database, &b) // Ignore errors for duplicates, etc.
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func parseTags(tags string) []string {
	if tags == "" {
		return nil
	}
	parts := strings.Split(tags, ",")
	var out []string
	for _, t := range parts {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
