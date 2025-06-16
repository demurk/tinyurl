package config

import (
	"flag"
	"fmt"
)

var (
	OriginURL     = flag.String("a", "localhost:8080", "Origin server url")
	BaseResultURL = flag.String("b", "", "Base result server url")
	ResultURL     = fmt.Sprintf("http://%s/%s", *OriginURL, *BaseResultURL)
)

func ParseFlags() {
	flag.Parse()
}
