package storage

import (
	"github.com/demurk/tinyurl/internal/config"
	"github.com/demurk/tinyurl/internal/db"
	"github.com/demurk/tinyurl/internal/types"
)

type StorageStruct struct {
	Get      func(string) (string, error)
	Set      func(string) (string, error)
	SetBatch func([]types.BatchJSONPostRequestData) ([]types.BatchJSONPostResponseData, error)
}

var storage StorageStruct

func Initialize() {
	conn := db.GetConnection()
	if conn != nil && conn.Ping() == nil {
		CreateTables()
		storage = StorageStruct{
			Get:      dbGetFullURL,
			Set:      dbSetFullURL,
			SetBatch: dbSetFullURLBatch,
		}
	} else if *config.FileStoragePath != "" {
		RestoreURLsFromFile()
		storage = StorageStruct{
			Get:      mGetFullURL,
			Set:      dSetFullURL,
			SetBatch: dSetFullURLBatch,
		}
	} else {
		storage = StorageStruct{
			Get:      mGetFullURL,
			Set:      mSetFullURL,
			SetBatch: mSetFullURLBatch,
		}
	}
}

func Get() StorageStruct {
	return storage
}
