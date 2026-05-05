package imgproc

import (
	"github.com/h2non/bimg"
	"os"
	"path/filepath"
	"fmt"
)

type Resizer struct{
	StoragePath string
}

func NewResizer(path string) (*Resizer, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	return &Resizer{StoragePath: path}, nil
}

func (r *Resizer) ProcessImage(imageBytes []byte,variantName string,baseName string)(string,error){
	mp := map[string]int{
		"thumb":200,
		"medium":800,
		"large":1600,
	}
	// set processing options
	options := bimg.Options{
		Width: mp[variantName],
		Type: bimg.WEBP,
		Quality: 80,
	}
	// resize the image
	resizedImage,err := bimg.Resize(imageBytes,options)
	if err != nil{
		return "",err
	}
	fileName := fmt.Sprintf("%s_%s.webp",baseName,variantName)
	imagePath := filepath.Join(r.StoragePath, fileName)
	err = os.WriteFile(imagePath,resizedImage,0644)
	if err != nil {
		return "", fmt.Errorf("write: %w", err)
	}
	return fileName,nil
}


