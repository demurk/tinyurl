package config

import "flag"

var (
	OriginURL = flag.String("a", "localhost:8080", "Origin server url")
	ResultURL = flag.String("b", "http://localhost:8080/", "Result server url")
)

func ParseFlags() {
	flag.Parse()
}
