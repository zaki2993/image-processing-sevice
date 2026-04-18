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
		"meduim":800,
		"large":1600,
	}
	options := bimg.Options{
		Width: mp[variantName],
		Type: bimg.WEBP,
		Quality: 80,
	}
	resizedImage,err := bimg.Resize(imageBytes,options)
	if err != nil{
		return "",err
	}
	filename := fmt.Sprintf("%s_%s.webp",baseName,variantName)
	imagePath := filepath.Join(r.StoragePath, filename)
	err = os.WriteFile(imagePath,resizedImage,0644)
	if err != nil {
		return "", fmt.Errorf("write: %w", err)
	}
	return filename,nil
}


