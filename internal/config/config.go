package config

import "os"
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

func Load() Configs{
	port := getEnv("PORT_GO","8081")
	storagepath := getEnv("STORAGE_PATH_GO","~/storage")
	return Configs{
		Port: port,
		StoragePath: storagepath,
	}
}
