package urls_storage

import (
	"github.com/demurk/tinyurl/internal/db"
)

type StorageStruct struct {
	Get func(string) (string, error)
	Set func(string) (string, error)
}

var storage StorageStruct

func New() {
	err := db.GetConnection().Ping()
	if err == nil {
		storage = StorageStruct{
			Get: dbGetFullURL,
			Set: dbSetFullURL,
		}
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
