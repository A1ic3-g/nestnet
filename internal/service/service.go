package service

import (
	"log"
	"net/http"
)

const ADDR = "0.0.0.0:8080"

// Test handler
func testHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello, world!\n"))
	if err != nil {
		log.Fatal(err)
	}
}

// Hello handler
//
// Send: hello message
//
// Receive: signed message
func helloHandler(w http.ResponseWriter, r *http.Request) {
}

// Posts handler
//
// Send: number of most recent posts to get (optional)
//
// Receive: posts with signature
func postsHandler(w http.ResponseWriter, r *http.Request) {

}

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", testHandler)

	err := http.ListenAndServe(ADDR, mux)
	if err != nil {
		log.Fatal(err)
	}
}
