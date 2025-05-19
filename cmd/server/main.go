// gomark/cmd/server/main.go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome to GoMark!")
    })

    fmt.Println("Server listening on :8080")
    http.ListenAndServe(":8080", nil)
}
