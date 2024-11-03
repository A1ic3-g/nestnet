package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
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

// ImageUploadRequest is the expected JSON structure for the POST request
type ImageUploadRequest struct {
	ImageBase64 string `json:"image_base64"` // Base64-encoded image data
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

const imageDir = "/var/lib/nestnet/images"

// imageHandler handles /image endpoint with GET and POST requests
func imageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Handle GET request
		md5Hash := r.URL.Query().Get("md5")
		if md5Hash == "" {
			http.Error(w, "md5 query parameter is required", http.StatusBadRequest)
			return
		}

		imagePath := filepath.Join(imageDir, md5Hash+".png")
		file, err := os.Open(imagePath)
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "Image not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to open image", http.StatusInternalServerError)
			}
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "image/png") // Set the appropriate content type
		if _, err := io.Copy(w, file); err != nil {
			http.Error(w, "Failed to serve image", http.StatusInternalServerError)
		}

	case http.MethodPost:
		// Handle POST request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Decode the Base64-encoded image data
		imageData, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			http.Error(w, "Failed to decode Base64 image", http.StatusBadRequest)
			log.Println("Error decoding Base64:", err)
			return
		}

		// Compute MD5 hash of the image data
		hash := md5.Sum(imageData)
		md5Hash := hex.EncodeToString(hash[:])

		// Save the image data as a PNG file
		imagePath := filepath.Join(imageDir, md5Hash+".png")
		file, err := os.Create(imagePath)
		if err != nil {
			http.Error(w, "Failed to save image file", http.StatusInternalServerError)
			log.Println("Error creating image file:", err)
			return
		}
		defer file.Close()

		// Write the image data to the file
		if _, err := file.Write(imageData); err != nil {
			http.Error(w, "Failed to write image file", http.StatusInternalServerError)
			log.Println("Error writing image data:", err)
			return
		}

		// Return the URL with the MD5 hash as a query parameter
		imageURL := fmt.Sprintf("/image?md5=%s", md5Hash)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(imageURL))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Start() {
	os.MkdirAll(imageDir, os.FileMode(0777))

	// Create a new ServeMux and register handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", testHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/posts", postsHandler)
	mux.HandleFunc("/image", imageHandler)

	// Start the server with ListenAndServe
	log.Printf("Server starting on %s\n", ADDR)
	if err := http.ListenAndServe(ADDR, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
