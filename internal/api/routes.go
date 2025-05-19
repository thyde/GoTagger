package api

import (
	"database/sql"
	"net/http"
)

// RegisterRoutes function to register the API routes for bookmarks to the HTTP server.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/api/bookmarks", ListBookmarksHandler(db))
	mux.HandleFunc("/api/bookmarks/add", AddBookmarkHandler(db))
	mux.HandleFunc("/api/bookmarks/search", SearchBookmarksHandler(db))
}
