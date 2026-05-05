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

// the function that handels route /resize 
func (h *Handler) Resize(w http.ResponseWriter,r *http.Request){
	// limit requests body to 10 MB
	r.Body = http.MaxBytesReader(w,r.Body,maxBytes)
	if err := r.ParseMultipartForm(maxBytes); err != nil{
		http.Error(w,"invalid upload",http.StatusBadRequest)
		return
	}
	// check if image field exist ,if it exists open the steam of data
	file,_,err := r.FormFile("image")
	if err != nil{
		http.Error(w,"image field is required",http.StatusBadRequest)
		return
	}
	defer file.Close()
	// read the steam data and store it in imageBytes
	imageBytes,err := io.ReadAll(file)
	if err != nil{
		http.Error(w,"can't read file",http.StatusInternalServerError)
		return
	}
	// check if image type is jpeg or png
	// read only 512 byte , it's enough to specify the image type
	readTo := 512
	if ln := len(imageBytes);ln < readTo{
		readTo = ln
	}
	contentType := http.DetectContentType(imageBytes[:readTo])
	if contentType != "image/jpeg" && contentType != "image/png" {
    http.Error(w, "invalid image type", http.StatusBadRequest)
    return
}
// generate unique a unique name for each image
unique := uuid.NewString()
// to images names
var thumb, medium, large string
var thumbErr, mediumErr, largeErr error
var wg sync.WaitGroup

// process three image sizes concurrently
wg.Add(3) 
go func() {
    defer wg.Done()  
    thumb, thumbErr = h.Resizer.ProcessImage(imageBytes, "thumb", unique)
}()

go func() {
    defer wg.Done()
    medium, mediumErr = h.Resizer.ProcessImage(imageBytes, "medium", unique)
}()

go func() {
    defer wg.Done()
    large, largeErr = h.Resizer.ProcessImage(imageBytes, "large", unique)
}()
// end of precessing
wg.Wait()  
if thumbErr != nil || mediumErr != nil || largeErr != nil{
	http.Error(w,"could not process image",http.StatusInternalServerError)
	return
}

// return the http response using json with paths of where each image is stored
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(resizeResponse{
    Thumb:  thumb,
    Medium: medium,
    Large:  large,
})
		}

