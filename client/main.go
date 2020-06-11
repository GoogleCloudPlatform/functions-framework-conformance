// This binary contains a valiation framework for functions frameworks.
package main

import (
	"flag"
	"log"
)

var (
	cmd          = flag.String("cmd", "", "command to run a Functions Framework server at localhost:8080")
	functionType = flag.String("type", "http", "type of function to validate (must be 'http', 'cloudevent', or 'legacyevent'")
)

func main() {
	flag.Parse()
	log.Printf("Validating %q for %s...", *cmd, *functionType)

	shutdown, err := start(*cmd)
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}

	if err := validate("http://localhost:8080", *functionType); err != nil {
		shutdown()
		log.Fatalf("Validation failure: %v", err)
	}

	log.Printf("All validation passed!")
	shutdown()
}
