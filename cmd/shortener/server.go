package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postPage)
	mux.HandleFunc(`/{id}`, idPage)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
