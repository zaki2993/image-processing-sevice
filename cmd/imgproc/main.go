package main

import (
	"log"
	"net/http"
	"github.com/zaki2993/image-processing-service/internal/config"
	"github.com/zaki2993/image-processing-service/internal/httpx"
	"github.com/zaki2993/image-processing-service/internal/imgproc"
)

func main(){
	confs := config.Load() 
	port := confs.Port 
	storagePath := confs.StoragePath
	resizer,err := imgproc.NewResizer(storagePath)
	if err != nil{
		log.Fatal(err)
	}
	handler := imgproc.NewHandler(resizer)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health",httpx.Health)
	mux.HandleFunc("POST /resize",handler.Resize)
	if err := http.ListenAndServe(":"+port,httpx.Recovery(httpx.Logging(mux))); err != nil {
		log.Fatal(err)
	}
}
