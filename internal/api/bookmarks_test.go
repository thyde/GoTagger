package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		title TEXT,
		tags TEXT,
		favorite BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return db
}

func TestAddAndListBookmarks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Add a bookmark
	body := `{"url":"https://golang.org","title":"Go","tags":["go","programming"],"favorite":true}`
	req := httptest.NewRequest("POST", "/api/bookmarks/add", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	AddBookmarkHandler(db)(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	// List bookmarks
	req = httptest.NewRequest("GET", "/api/bookmarks", nil)
	w = httptest.NewRecorder()
	ListBookmarksHandler(db)(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var bookmarks []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &bookmarks); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(bookmarks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(bookmarks))
	}
}

func TestSearchBookmarks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, err := db.Exec(`INSERT INTO bookmarks (url, title, tags, favorite) VALUES (?, ?, ?, ?)`,
		"https://golang.org", "Go", "go,programming", true)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/bookmarks/search?q=go", nil)
	w := httptest.NewRecorder()
	SearchBookmarksHandler(db)(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var bookmarks []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &bookmarks); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(bookmarks) == 0 {
		t.Fatalf("expected at least 1 bookmark, got %d", len(bookmarks))
	}
}
