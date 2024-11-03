package service

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"nestnet/internal/database"
	"nestnet/internal/database/generated"
	"net/http"
	"os"
	"path/filepath"
)

// ImageUploadRequest is the expected JSON structure for the POST request
type ImageUploadRequest struct {
	ImageBase64 string `json:"image_base64"` // Base64-encoded image data
}

const ADDR = ":8080"

// defaultHandler gives a hello world message as a default response
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello, world!\n"
	if _, err := w.Write([]byte(msg)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Println("Error writing response:", err)
	}
}

// postsHandler sends the most recent posts as JSON
func postsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received posts request")
	var posts []generated.Post
	testPost := generated.Post{
		ID:     uuid.New().String(),
		Title:  "Test title",
		Body:   "lorem ipsum yeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee yippeeeee",
		Imgmd5: "",
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

// addPostHandler handles adding a post to the user's posts
func addPostHandler(w http.ResponseWriter, r *http.Request) {
	var post generated.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}

	database.AddPost(post)
	w.WriteHeader(http.StatusCreated)
}

// addPeerHandler handles adding a peer to the user's peers
func addPeerHandler(w http.ResponseWriter, r *http.Request) {
	var peer generated.Peer
	err := json.NewDecoder(r.Body).Decode(&peer)
	if err != nil {
		log.Fatal(err)
	}

	database.AddPeer(peer)
	w.WriteHeader(http.StatusCreated)
}

func Start() {
	os.MkdirAll(imageDir, os.FileMode(0777))

	// Create a new ServeMux and register handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)
	mux.HandleFunc("/posts", postsHandler)
	mux.HandleFunc("/image", imageHandler)
	mux.HandleFunc("/add_post", addPostHandler)
	mux.HandleFunc("/add_peer", addPeerHandler)

	// Start the server with ListenAndServe
	log.Printf("Server starting on %s\n", ADDR)
	if err := http.ListenAndServe(ADDR, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
