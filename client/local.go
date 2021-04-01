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
	"syscall"
	"time"
)

type localFunctionServer struct {
	output string
	cmd    string
}

func (l *localFunctionServer) Start() (func(), error) {
	args := strings.Fields(l.cmd)
	cmd := exec.Command(args[0], args[1:]...)

	// Set a process group ID so that later we can kill child processes too. As an
	// example, if the command is `go run main.go`, Go will build a binary in a
	// temp dir and then execute it. If we simply cmd.Process.Kill() the exec.Command
	// then the running binary will not be killed. Only if we make a group and then
	// kill the group will child processes be killed.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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

	// Give it some time to do its setup.
	time.Sleep(time.Duration(*startDelay) * time.Second)

	shutdown := func() {
		stdout.Close()
		stderr.Close()

		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err != nil {
			log.Printf("Failed to get pgid: %v", err)

			// Kill just the parent process since we failed to get the process group ID.
			if err := cmd.Process.Kill(); err != nil {
				log.Fatalf("Failed to kill process: %v", err)
			}
		} else {
			// Kill the whole process group.
			if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
				log.Fatalf("Failed to kill process group: %v", err)
			}
		}

		log.Printf("Framework server shut down. Wrote logs to %v and %v.", stdoutFile, stderrFile)
	}
	return shutdown, nil
}

func (l *localFunctionServer) OutputFile() ([]byte, error) {
	return ioutil.ReadFile(l.output)
}
