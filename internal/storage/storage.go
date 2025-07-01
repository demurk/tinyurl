package storage

import "github.com/demurk/tinyurl/internal/config"

type StorageStruct struct {
	Get func(string) (string, error)
	Set func(string) (string, error)
}

var storage StorageStruct

func Initialize() {
	if *config.FileStoragePath != "" {
		storage = StorageStruct{
			Get: mGetFullURL,
			Set: dSetFullURL,
		}
		RestoreURLsFromFile()
	} else {
		storage = StorageStruct{
			Get: mGetFullURL,
			Set: mSetFullURL,
		}
	}
}

func Get() StorageStruct {
	return storage
}
