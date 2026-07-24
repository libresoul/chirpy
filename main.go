package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	server := http.Server{Addr: ":8080"}
	fmt.Println("Server started on port 8080")
	server.ListenAndServe()
}
