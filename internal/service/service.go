package service

import (
	"log"
	"net/http"
)

const ADDR = "0.0.0.0:8080"

func testHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello, world!\n"))
	if err != nil {
		log.Fatal(err)
	}
}

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", testHandler)

	err := http.ListenAndServe(ADDR, mux)
	if err != nil {
		log.Fatal(err)
	}
}
