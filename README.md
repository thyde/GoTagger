# GoMark

GoMark is a modern, web-based bookmark manager written in Go. It allows users to save, tag, search, edit, and favorite bookmarks through a beautiful web interface and a REST API, with persistent storage and optional metadata enrichment.

## Features

- Add, edit, delete, and list bookmarks
- Tagging and keyword search
- Favorite bookmarks for quick access
- Pagination and sorting
- Tag cloud with filtering and tag usage counts
- Import/export bookmarks (JSON)
- Duplicate detection on add
- Modern, responsive web UI (Google-like aesthetic)
- Persistent storage with SQLite (default) or PostgreSQL (optional)
- JSON-based REST API
- Written with idiomatic, testable Go

## Getting Started

### Prerequisites

- Go 1.21+
- SQLite (default, install via Homebrew on macOS: `brew install sqlite`)
- (Optional) PostgreSQL for advanced setups

### Installation & Running

```bash
git clone https://github.com/thyde/GoTagger.git
cd GoTagger
go run cmd/server/main.go
```

The server will start on [http://localhost:8080](http://localhost:8080).

### Usage

- Visit [http://localhost:8080](http://localhost:8080) in your browser.
- Use the web interface to add, edit, delete, search, and favorite bookmarks.
- Filter bookmarks by tag using the tag cloud.
- Import/export bookmarks as JSON.
- The database is always seeded to include at least 5 bookmarks (Google, Bing, Yahoo, Amazon, and examples).
- Bookmarks are deduplicated by URL on startup.

### API Endpoints

- `GET /api/bookmarks` — List bookmarks (supports `search`, `page`, `per_page`, `sort` query params)
- `POST /api/bookmarks/add` — Add a new bookmark (JSON body)
- `POST /api/bookmarks/update` — Update a bookmark (JSON body)
- `GET /api/bookmarks/delete?id=...` — Delete a bookmark by ID
- `GET /api/bookmarks/search?q=...` — Search bookmarks by tag, title, or URL
- `GET /api/tags` — List all tags and their usage counts
- `GET /api/bookmarks/export` — Export bookmarks as JSON
- `POST /api/bookmarks/import` — Import bookmarks from JSON

### Testing

Run tests with:

```bash
go test ./internal/api
```

## Project Structure

- `cmd/server/` — Main server entrypoint and static files
- `internal/api/` — API handlers and routes
- `internal/db/` — Database setup and logic
- `internal/model/` — Data models
- `internal/service/` — Background services (e.g. metadata)
- `test/` — API and integration tests

## License

MIT
