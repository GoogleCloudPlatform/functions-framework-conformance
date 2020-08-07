// Copyright 2020 Google LLC
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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type localFunctionServer struct {
	output string
	cmd    string
}

func (l *localFunctionServer) Start() (func(), error) {
	args := strings.Fields(l.cmd)
	cmd := exec.Command(args[0], args[1:]...)

	stdout, err := os.Create(stdoutFile)
	if err != nil {
		return nil, err
	}
	cmd.Stdout = stdout

	stderr, err := os.Create(stderrFile)
	if err != nil {
		return nil, err
	}
	cmd.Stderr = stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	log.Printf("Framework server started.")

	// Give it a second to do its setup.
	time.Sleep(time.Second)

	shutdown := func() {
		// TODO: kill processes properly.
		if err := cmd.Process.Kill(); err != nil {
			log.Fatalf("failed to kill process: %v", err)
		}
		stdout.Close()
		stderr.Close()
		log.Printf("Framework server shut down. Wrote logs to %v and %v.", stdoutFile, stderrFile)
	}
	return shutdown, nil
}

func (l *localFunctionServer) OutputFile() ([]byte, error) {
	return ioutil.ReadFile(l.output)
}
