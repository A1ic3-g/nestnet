package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type HelloReq struct {
	Addr string   `json:"addr"`
	X    *big.Int `json:"x"` // Public key x
	Y    *big.Int `json:"y"` // Public key y
}

type HelloRes struct {
	R *big.Int `json:"r"` // Signature r
	S *big.Int `json:"s"` // Signature s
}

type Post struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	ImgMd5  string    `json:"img_md5"`
	ImgName string    `json:"img_name"`
}

const ADDR = ":8080"

// Test handler
func testHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello, world!\n"
	if _, err := w.Write([]byte(msg)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Println("Error writing response:", err)
	}
}

// helloHandler asks the given address to sign a hello message and verifies it using the given public key
func helloHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody HelloReq
	msg := "HELLO"

	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	out, err := http.NewRequest(http.MethodPost, reqBody.Addr, strings.NewReader(msg))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Println("Error creating request:", err)
		return
	}

	res, err := http.DefaultClient.Do(out)
	if err != nil {
		http.Error(w, "Request failed", http.StatusInternalServerError)
		log.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	var resBody HelloRes
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		log.Println("Error decoding response:", err)
		return
	}

	hash := sha256.Sum256([]byte(msg))
	success := ecdsa.Verify(&ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     reqBody.X,
		Y:     reqBody.Y,
	}, hash[:], resBody.R, resBody.S)

	successStr := strconv.FormatBool(success)
	if _, err := w.Write([]byte(successStr)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Println("Error writing response:", err)
	}
}

// postsHandler sends the most recent posts as JSON
func postsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received posts request")
	var posts []Post
	testPost := Post{
		ID:      uuid.New(),
		Title:   "Test title",
		Body:    "lorem ipsum yeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee yippeeeee",
		ImgMd5:  "",
		ImgName: "",
	}

	posts = append(posts, testPost)

	// Marshal the posts slice into JSON
	jsonData, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		log.Println("Error marshalling JSON:", err)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON data to the response
	if _, err := w.Write(jsonData); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Println("Error writing response:", err)
	}
}

func Start() {
	// Create a new ServeMux and register handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", testHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/posts", postsHandler)

	// Start the server with ListenAndServe
	log.Printf("Server starting on %s\n", ADDR)
	if err := http.ListenAndServe(ADDR, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
