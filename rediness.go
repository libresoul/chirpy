package main

import (
	"log"
	"net/http"
)

func handlerRediness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatal("Failed to write response: ", err)
	}
}
