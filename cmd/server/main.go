package main

import (
    "net/http"

    "github.com/traceylum1/observability-api/internal/server"
)

func main() {
    r := server.NewRouter()
    http.ListenAndServe(":3000", r)
}