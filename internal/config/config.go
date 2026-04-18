package config

import(
	"log"
	"os"
)
type Configs struct{
	Port string
	StoragePath string
}

func getEnv(key ,fullback string) string{
	if val := os.Getenv(key); val != ""{
		return val
	}
	return fullback
}

func Load() Configs {
    storagePath := os.Getenv("IMGPROC_STORAGE_PATH")
    if storagePath == "" {
        log.Fatal("IMGPROC_STORAGE_PATH must be set")
    }
    return Configs{
        Port:        getEnv("IMGPROC_PORT", "8081"),
        StoragePath: storagePath,
    }
}
