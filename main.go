// A web service that checks the build secret its image was built with
// against an expected value, without ever showing the secret or its hash.
// The build writes only the secret's SHA-256 into the image; at runtime the
// service compares it with EXPECTED_SHA256 and answers matched or not
// matched.
package main

import (
	"crypto/subtle"
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
	built := strings.TrimSpace(string(sum))
	expected := strings.ToLower(strings.TrimSpace(os.Getenv("EXPECTED_SHA256")))

	result := "not matched"
	if expected != "" && subtle.ConstantTimeCompare([]byte(built), []byte(expected)) == 1 {
		result = "matched"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("Build secret check: " + result + "\n"))
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
