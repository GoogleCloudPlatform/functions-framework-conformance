// This binary starts an HTTP server to serve the Go FF validation test functions.
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// HTTP is a simple HTTP function that writes the request body to the response body.
func HTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := ioutil.WriteFile("function_output.json", body, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", HTTP)

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
