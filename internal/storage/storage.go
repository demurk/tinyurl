package storage

import (
	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/db"
	"github.com/demurk/tinyurl/internal/types"
)

type StorageStruct struct {
	Get      func(string) (string, error)
	Set      func(string) (string, error)
	SetBatch func([]types.BatchJsonPostRequestData) ([]types.BatchJsonPostResponseData, error)
}

var storage StorageStruct

func Initialize() {
	err := db.GetConnection().Ping()
	if err == nil {
		storage = StorageStruct{
			Get:      dbGetFullURL,
			Set:      dbSetFullURL,
			SetBatch: dbSetFullURLBatch,
		}
	} else if *config.FileStoragePath != "" {
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
