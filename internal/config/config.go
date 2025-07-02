package config

import (
	"flag"
	"os"
)

var (
	OriginURL       *string
	ResultURL       *string
	DBConnectionURL *string
	FileStoragePath *string
)

func Parse() {
	OriginURL = flag.String("a", "localhost:8080", "Origin server url")
	originURLEnv, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		OriginURL = &originURLEnv
	}

	ResultURL = flag.String("b", "http://localhost:8080", "Result server url")
	resultURLEnv, exists := os.LookupEnv("BASE_URL")
	if exists {
		ResultURL = &resultURLEnv
	}

	DBConnectionURL = flag.String("d", "host='127.0.0.1' port='5432' user='postgres' password='password' dbname='tinyurl' sslmode=disable", "Database connection URL")
	dbConnectionURLEnv, exists := os.LookupEnv("DATABASE_DSN")
	if exists {
		DBConnectionURL = &dbConnectionURLEnv
	}

	FileStoragePath = flag.String("f", "./urls_storage.jsonl", "Urls storage file path (JSONL format)")
	fileStoragePathEnv, exists := os.LookupEnv("FILE_STORAGE_PATH")
	if exists {
		FileStoragePath = &fileStoragePathEnv
	}

	flag.Parse()
}
