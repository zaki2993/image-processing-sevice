package main

import (
	"encoding/json"
	"log"
	"net/http"
)
func health(w http.ResponseWriter, r *http.Request){
	response := map[string]string{
		"status":"ok",
	}
	w.Header().Set("Content-Type","application/json")
json.NewEncoder(w).Encode(response)
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health",health)
	err := http.ListenAndServe(":8081",mux)
	if err != nil{
		log.Fatal(err)
	}

}
