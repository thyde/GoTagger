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
	mux.HandleFunc("/api/bookmarks/update", UpdateBookmarkHandler(db))
	mux.HandleFunc("/api/bookmarks/delete", DeleteBookmarkHandler(db))
	mux.HandleFunc("/api/tags", ListTagsHandler(db))
	mux.HandleFunc("/api/bookmarks/export", ExportBookmarksHandler(db))
	mux.HandleFunc("/api/bookmarks/import", ImportBookmarksHandler(db))
}
