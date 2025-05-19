package db

import (
	"GoTagger/internal/model"
	"database/sql"
	"fmt"
	"strings"
)

// CreateBookmark inserts a new bookmark into the database.
func CreateBookmark(db *sql.DB, b *model.Bookmark) error {
	tags := strings.Join(b.Tags, ",")
	res, err := db.Exec(`INSERT INTO bookmarks (url, title, tags, favorite) VALUES (?, ?, ?, ?)`, b.URL, b.Title, tags, b.Favorite)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		b.ID = id
	}
	return err
}

// GetBookmarkByID retrieves a bookmark by its ID.
func GetBookmarkByID(db *sql.DB, id int64) (*model.Bookmark, error) {
	var b model.Bookmark
	var tags string
	err := db.QueryRow(`SELECT id, url, title, tags, favorite, created_at, updated_at FROM bookmarks WHERE id = ?`, id).
		Scan(&b.ID, &b.URL, &b.Title, &tags, &b.Favorite, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	b.Tags = parseTags(tags)
	return &b, nil
}

// ListBookmarks returns all bookmarks in the database.
func ListBookmarks(db *sql.DB) ([]model.Bookmark, error) {
	rows, err := db.Query(`SELECT id, url, title, tags, favorite, created_at, updated_at FROM bookmarks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bookmarks []model.Bookmark
	for rows.Next() {
		var b model.Bookmark
		var tags string
		if err := rows.Scan(&b.ID, &b.URL, &b.Title, &tags, &b.Favorite, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		b.Tags = parseTags(tags)
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

// UpdateBookmark updates an existing bookmark by ID.
func UpdateBookmark(db *sql.DB, b *model.Bookmark) error {
	tags := strings.Join(b.Tags, ",")
	_, err := db.Exec(`UPDATE bookmarks SET url = ?, title = ?, tags = ?, favorite = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, b.URL, b.Title, tags, b.Favorite, b.ID)
	return err
}

// DeleteBookmark deletes a bookmark by ID.
func DeleteBookmark(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM bookmarks WHERE id = ?`, id)
	return err
}

// SearchBookmarksByKeywordOrTag searches bookmarks by keyword in tags, title, or url.
func SearchBookmarksByKeywordOrTag(db *sql.DB, q string) ([]model.Bookmark, error) {
	rows, err := db.Query(`SELECT id, url, title, tags, favorite, created_at, updated_at FROM bookmarks WHERE tags LIKE ? OR title LIKE ? OR url LIKE ?`, "%"+q+"%", "%"+q+"%", "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bookmarks []model.Bookmark
	for rows.Next() {
		var b model.Bookmark
		var tags string
		if err := rows.Scan(&b.ID, &b.URL, &b.Title, &tags, &b.Favorite, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		b.Tags = parseTags(tags)
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

// ListBookmarksPaginated returns bookmarks with pagination, search, and sorting.
func ListBookmarksPaginated(db *sql.DB, search, sort string, page, perPage int) ([]model.Bookmark, int, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = "WHERE tags LIKE ? OR title LIKE ? OR url LIKE ?"
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	order := "ORDER BY created_at DESC"
	if sort == "title" {
		order = "ORDER BY title COLLATE NOCASE ASC"
	} else if sort == "favorite" {
		order = "ORDER BY favorite DESC, created_at DESC"
	}
	limit := "LIMIT ? OFFSET ?"
	args = append(args, perPage, (page-1)*perPage)

	query := fmt.Sprintf(`SELECT id, url, title, tags, favorite, created_at, updated_at FROM bookmarks %s %s %s`, where, order, limit)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var bookmarks []model.Bookmark
	for rows.Next() {
		var b model.Bookmark
		var tags string
		if err := rows.Scan(&b.ID, &b.URL, &b.Title, &tags, &b.Favorite, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, 0, err
		}
		b.Tags = parseTags(tags)
		bookmarks = append(bookmarks, b)
	}
	// Get total count
	total := 0
	countQuery := "SELECT COUNT(*) FROM bookmarks " + where
	row := db.QueryRow(countQuery, args[:len(args)-2]...)
	row.Scan(&total)
	return bookmarks, total, nil
}

// parseTags splits a comma-separated tag string into a slice.
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
