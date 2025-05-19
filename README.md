# GoMark

GoMark is a lightweight, web-based bookmark manager written in Go. It allows users to save, tag, and search bookmarks through a simple REST API, with persistent storage and optional metadata enrichment.

## Features

- Add, edit, delete, and list bookmarks
- Tagging and keyword search
- Favorite bookmarks for quick access
- Background metadata fetching (e.g. page title)
- Persistent storage with SQLite or PostgreSQL
- JSON-based REST API
- Written with idiomatic, testable Go

## Getting Started

### Prerequisites

- Go 1.21+
- SQLite (default) or PostgreSQL (optional)

### Install

```bash
git clone https://github.com/your-username/gomark.git
cd gomark
go run main.go
