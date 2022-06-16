// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const testProgram = `package main

import (
	"fmt"
	"time"
)

func main() {
	sleepDuration := time.Second * 90
	fmt.Printf("Hello from test program. Sleeping for %v.\n", sleepDuration)
	time.Sleep(sleepDuration)
}
`

func TestStartAndShutdown(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "main.go")
	if err := ioutil.WriteFile(f, []byte(testProgram), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	outputFile := filepath.Join(dir, "function_output.json")

	server := localFunctionServer{
		// `go run` compiles the program and then execs it, which allows us to test that
		// the whole process group is killed. It's done this way instead of something
		// simpler like "/bin/sh -c 'exec sleep 90'" so that it's cross-platform compatible.
		cmd: fmt.Sprintf("go run %s", f),
	}

	shutdown, err := server.Start(defaultStdoutFile, defaultStderrFile, outputFile)
	if shutdown != nil {
		defer shutdown()
	}

	if err != nil {
		t.Errorf("unable to start localFunctionServer: %v", err)
	}
}
