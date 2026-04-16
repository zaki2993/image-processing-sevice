package imgproc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const maxBytes = 10 << 20

type resizeResponse struct {
	Path string `json:"path"`
}

type Handler struct{
	Resizer *Resizer
}

func NewHandler(r *Resizer) *Handler {
	return &Handler{Resizer: r}
}

func (h *Handler) Resize(w http.ResponseWriter,r *http.Request){
	r.Body = http.MaxBytesReader(w,r.Body,maxBytes)
	if err := r.ParseMultipartForm(maxBytes); err != nil{
		http.Error(w,"invalid upload",http.StatusBadRequest)
		return
	}
	file,_,err := r.FormFile("image")
	if err != nil{
		http.Error(w,"image field is required",http.StatusBadRequest)
		return
	}
	defer file.Close()
	imageBytes,err := io.ReadAll(file)
	if err != nil{
		http.Error(w,"can't read file",http.StatusInternalServerError)
		return
	}
	filename,err := h.Resizer.ProcessImage(imageBytes,uuid.NewString())
	if err != nil{
		http.Error(w,fmt.Sprint(err),http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(resizeResponse{Path: filename})
}

