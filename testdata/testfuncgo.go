// Package function is a Go test function.
package function

import (
	"fmt"
	"io/ioutil"
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
	if curr_dir, err := os.Getwd(); err != nil {
		fmt.Printf("Failed to get working directory: %s", err)
	} else {
		fmt.Printf("Writing output json to: %s", curr_dir)
	}
	if err := ioutil.WriteFile("function_output.json", body, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
