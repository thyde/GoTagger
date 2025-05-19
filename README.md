# GoMark

GoMark is a lightweight, web-based bookmark manager written in Go. It allows users to save, tag, search, and favorite bookmarks through a simple web interface and REST API, with persistent storage and optional metadata enrichment.

## Features

- Add, edit, delete, and list bookmarks
- Tagging and keyword search
- Favorite bookmarks for quick access
- Background metadata fetching (e.g. page title)
- Persistent storage with SQLite (default) or PostgreSQL (optional)
- JSON-based REST API
- Modern, responsive web UI (Google-like aesthetic)
- Written with idiomatic, testable Go

## Getting Started

### Prerequisites

- Go 1.21+
- SQLite (default, installed via Homebrew on macOS: `brew install sqlite`)
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
- Use the web interface to add, search, and favorite bookmarks.
- The database is automatically seeded to always include at least 5 bookmarks, including Google, Bing, Yahoo, and Amazon.
- Bookmarks are deduplicated by URL on startup.

### API Endpoints

- `GET /api/bookmarks` — List all bookmarks
- `POST /api/bookmarks/add` — Add a new bookmark (JSON body)
- `GET /api/bookmarks/search?q=term` — Search bookmarks by tag, title, or URL

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
