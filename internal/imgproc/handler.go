package imgproc

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

const maxBytes = 10 << 20

type resizeResponse struct {
    Thumb  string `json:"thumb"`
    Medium string `json:"medium"`
    Large  string `json:"large"`
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
unique := uuid.NewString()
var thumb, medium, large string
var thumbErr, mediumErr, largeErr error
var wg sync.WaitGroup

wg.Add(3) 

go func() {
    defer wg.Done()  
    thumb, thumbErr = h.Resizer.ProcessImage(imageBytes, 200, unique)
}()

go func() {
    defer wg.Done()
    medium, mediumErr = h.Resizer.ProcessImage(imageBytes, 800, unique)
}()

go func() {
    defer wg.Done()
    large, largeErr = h.Resizer.ProcessImage(imageBytes, 1600, unique)
}()

wg.Wait()  
if thumbErr != nil || mediumErr != nil || largeErr != nil{
	http.Error(w,"could not process image",http.StatusInternalServerError)
	return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(resizeResponse{
    Thumb:  thumb,
    Medium: medium,
    Large:  large,
})
		}

