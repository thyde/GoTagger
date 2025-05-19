package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "GoTagger.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		title TEXT,
		tags TEXT,
		favorite BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)
	return err
}

func SeedDummyData(db *sql.DB) error {
	_, err := db.Exec(`
	INSERT INTO bookmarks (url, title, tags, favorite) VALUES
	('https://golang.org', 'The Go Programming Language', 'go,programming,language', 1),
	('https://github.com', 'GitHub', 'code,repository,git', 0),
	('https://news.ycombinator.com', 'Hacker News', 'news,tech,startups', 0),
	('https://reddit.com/r/golang', 'Reddit Go', 'go,community,forum', 1)
	;`)
	return err
}
