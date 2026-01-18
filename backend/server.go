package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func main() {
	// http.HandleFunc("/posts", postsHandler)
	http.Handle("/posts", enableCORS(http.HandlerFunc(postsHandler)))

	fmt.Println("Server is running at http:localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
