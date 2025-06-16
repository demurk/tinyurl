package config

import (
	"flag"
	"fmt"
)

var (
	OriginURL = flag.String("a", "localhost:8080", "Origin server url")
	ResultURL = fmt.Sprintf("%s/%s", *OriginURL, *flag.String("b", "", "Base result server url"))
)

func ParseFlags() {
	flag.Parse()
}
