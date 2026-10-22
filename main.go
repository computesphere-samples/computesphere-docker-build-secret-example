// A web service that reports the SHA-256 of the build secret its image was
// built with. The secret itself is never stored in the image.
package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	sum, err := os.ReadFile("/secret.sha256")
	if err != nil {
		log.Fatalf("read the build secret's hash: %v", err)
	}
	hash := strings.TrimSpace(string(sum))
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("build secret sha256: " + hash + "\n"))
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
