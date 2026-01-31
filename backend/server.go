package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
)

type Post struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	Email   string `json:"email"`
	Domain  string `json:"domain"`
	// ID      int    `json:"id"`
	// Body    string `json:"body"`
	// Title   string `json:"title"`
	// Content string `json:"content"`
}

var (
	posts   = make(map[int]Post)
	nextID  = 1
	postsMu sync.Mutex
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// handler serves the index.html file
func handler(w http.ResponseWriter, r *http.Request) {
	// Open the HTML file
	file, err := os.Open("./frontend/index.html")
	if err != nil {
		// If file not found, return 404
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close() // Ensure the file is closed when done

	// Get file info to set the Content-Length header
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error getting file info: %v", err)
		return
	}

	// Set headers (optional but good practice)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Write the file content to the response
	http.ServeContent(w, r, "index.html", fileInfo.ModTime(), file)
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin": //macOS
		cmd = "open"
	default: //Linux, BSD, etc.
		cmd = "xdg-open"
	}

	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}

func main() {
	// http.HandleFunc("/posts", postsHandler)
	http.Handle("/posts", enableCORS(http.HandlerFunc(postsHandler)))

	fmt.Println("Server is running at http:localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	http.HandleFunc("/", handler)
	addr := ":3000"
	url := "http://localhost" + addr

	//Start the server
	go func() {
		log.Printf("Starting the server on %s", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Fatalf("Listen and serve error: %v", err)
		}
	}()

	//Open the browser after the server is up
	log.Printf("Opening the server at: %s in your default browser", url)
	err := openBrowser(url)
	if err != nil {
		log.Printf("The server failed to open at %v, visit %v yourself", err, url)
	}

}

// PostsHandler
func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetPosts(w, r)
	case "POST":
		handlePostPosts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET POSTS
func handleGetPosts(w http.ResponseWriter, _ *http.Request) {
	postsMu.Lock() //Locks the goroutine to avoid race conditions

	defer postsMu.Unlock() //Unlocks the goroutine at the end

	ps := make([]Post, 0, len(posts))
	for _, p := range posts {
		ps = append(ps, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ps)
}

// POST POSTS
func handlePostPosts(w http.ResponseWriter, r *http.Request) {
	var p Post

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &p); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Posted data: %+v\n", p)

	postsMu.Lock()
	defer postsMu.Unlock()

	// p.ID = nextID
	// nextID++
	// posts[p.ID] = p

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// // GET Post
// func handleGetPost(w http.ResponseWriter, _ *http.Request, id int) {
// 	postsMu.Lock()
// 	defer postsMu.Unlock()

// 	p, ok := posts[id]
// 	if !ok {
// 		http.Error(w, "Post not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(p)
// }

// // DELETE Post
// func handleDeletePost(w http.ResponseWriter, _ *http.Request, id int) {
// 	postsMu.Lock()
// 	defer postsMu.Unlock()

// 	_, ok := posts[id]
// 	if !ok {
// 		http.Error(w, "Post not found", http.StatusNotFound)
// 		return
// 	}

// 	delete(posts, id)
// 	w.WriteHeader(http.StatusOK)
// }
