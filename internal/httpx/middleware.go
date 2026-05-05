package httpx

import (
	"log"
	"time"
	"net/http"
)
// meddleware that handels time for all handlers
func Logging(next http.Handler) http.Handler{
	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)    
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })

}

// meddleware that handels unexpected errors for all handlers 
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v", err)
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
