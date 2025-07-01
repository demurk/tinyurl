package config

import (
	"flag"
	"os"
)

var (
	OriginURL *string
	ResultURL *string
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

	flag.Parse()
}
