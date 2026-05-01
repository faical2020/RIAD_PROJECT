package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve static files from the "dist" directory
	fs := http.FileServer(http.Dir("./dist"))
	
	// Handle all requests by serving the dist folder
	// To support SPA routing, we can add a custom handler that serves index.html for 404s
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if file exists, if not serve index.html (for SPA)
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		
		// Check if the file exists in the dist directory
		if _, err := os.Stat("./dist" + path); os.IsNotExist(err) {
			http.ServeFile(w, r, "./dist/index.html")
			return
		}
		
		fs.ServeHTTP(w, r)
	})

	log.Printf("Web server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
