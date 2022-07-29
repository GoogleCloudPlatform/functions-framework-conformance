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
	"strings"
	"time"
)

type localFunctionServer struct {
	functionOutputFile string
	cmd                string
	stdoutFile         string
	stderrFile         string
	envs			   string
}

func (l *localFunctionServer) Start(stdoutFile, stderrFile, functionOutputFile string) (func(), error) {
	l.stdoutFile = stdoutFile
	l.stderrFile = stderrFile
	l.functionOutputFile = functionOutputFile
	args := strings.Fields(l.cmd)
	cmd := newCmd(args)

	stdout, err := os.Create(l.stdoutFile)
	if err != nil {
		return nil, err
	}
	cmd.Stdout = stdout

	stderr, err := os.Create(l.stderrFile)
	if err != nil {
		return nil, err
	}
	cmd.Stderr = stderr
	cmd.Env = os.Environ()
	var variables = strings.Split(l.envs,",")
	for _,s := range variables{
		if s != "" {
			cmd.Env = append(cmd.Env, s)
		}
	}
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

		if err := stopCmd(cmd); err != nil {
			log.Fatalf("Failed to shut down framework server: %v", err)
		}

		log.Printf("Framework server shut down. Wrote logs to %v and %v.", l.stdoutFile, l.stderrFile)
	}
	return shutdown, nil
}

func (l *localFunctionServer) OutputFile() ([]byte, error) {
	return ioutil.ReadFile(l.functionOutputFile)
}
