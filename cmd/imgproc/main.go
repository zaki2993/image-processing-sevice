package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/zaki2993/image-processing-service/internal/config"
)
func health(w http.ResponseWriter, r *http.Request){
	response := map[string]string{
		"status":"ok",
	}
	w.Header().Set("Content-Type","application/json")
json.NewEncoder(w).Encode(response)
}

func main(){
	confs := config.Load()
	port := confs.Port
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health",health)
	err := http.ListenAndServe(":"+port,mux)
	if err != nil{
		log.Fatal(err)
	}

}
