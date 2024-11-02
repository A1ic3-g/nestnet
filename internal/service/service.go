package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type helloReq struct {
	addr string
	x    *big.Int // Public key x
	y    *big.Int // Public key y
}

type helloRes struct {
	r *big.Int // Signature r
	s *big.Int // Signature s
}

const ADDR = "0.0.0.0:8080"

// Test handler
func testHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello, world!\n"
	sent, err := w.Write([]byte(msg))
	if sent != len([]byte(msg)) || err != nil {
		log.Fatal(err)
	}
}

// helloHandler asks the given address to sign a hello message and verifies it using the given public key
func helloHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody helloReq
	msg := "HELLO"

	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		log.Fatal(err)
	}

	out, err := http.NewRequest(http.MethodPost, reqBody.addr, strings.NewReader(msg))
	res, err := http.DefaultClient.Do(out)
	if err != nil {
		log.Fatal(err)
	}

	var resBody helloRes
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		log.Fatal(err)
	}

	hash := sha256.Sum256([]byte(msg))
	success := ecdsa.Verify(&ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     reqBody.x,
		Y:     reqBody.y,
	}, hash[:], resBody.r, resBody.s)
	successStr := strconv.FormatBool(success)

	sent, err := w.Write([]byte(successStr))
	if sent != len([]byte(successStr)) || err != nil {
		log.Fatal(err)
	}
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
	mux.HandleFunc("/hello", helloHandler)

	err := http.ListenAndServe(ADDR, mux)
	if err != nil {
		log.Fatal(err)
	}
}
