package httpx

import (
    "encoding/json"
    "net/http"
)

// function that handels /health route
func Health(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
