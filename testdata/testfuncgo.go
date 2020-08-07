// Package function is a Go test function.
package function

import (
	"io/ioutil"
	"net/http"
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
