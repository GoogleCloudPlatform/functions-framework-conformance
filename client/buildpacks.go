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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	pack "github.com/buildpacks/pack/pkg/client"
	"github.com/buildpacks/pack/pkg/logging"
)

const (
	image             = "conformance-test-func"
	builderURL        = "gcr.io/buildpacks/builder:%s"
	gcfTargetPlatform = "gcf"
)

type buildpacksFunctionServer struct {
	functionOutputFile string
	source             string
	target             string
	funcType           string
	runtime            string
	tag                string
	ctID               string
	logStdout          *os.File
	logStderr          *os.File
	stdoutFile         string
	stderrFile         string
	envs               []string
}

func (b *buildpacksFunctionServer) Start(stdoutFile, stderrFile, functionOutputFile string) (func(), error) {
	b.functionOutputFile = functionOutputFile
	b.stdoutFile = stdoutFile
	b.stderrFile = stderrFile
	typ := *functionType
	if typ == "legacyevent" {
		typ = "event"
	}

	ctx := context.Background()
	if err := b.build(ctx); err != nil {
		return nil, fmt.Errorf("building function container: %v", err)
	}

	shutdown, err := b.run()
	if err != nil {
		return nil, fmt.Errorf("running function container: %v", err)
	}

	return shutdown, nil
}

func (b *buildpacksFunctionServer) OutputFile() ([]byte, error) {
	cmd := exec.Command("docker", "cp", filepath.Join(fmt.Sprintf("%s:/workspace", b.containerID()), b.functionOutputFile), os.TempDir())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to copy output file from the container: %v: %s", err, string(output))
	}
	return ioutil.ReadFile(filepath.Join(os.TempDir(), filepath.Base(b.functionOutputFile)))
}

func (b *buildpacksFunctionServer) build(ctx context.Context) error {
	builder := fmt.Sprintf(builderURL, b.tag)

	cmd := exec.Command("docker", "pull", builder)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull builder image %s: %v: %s", builder, err, string(output))
	}

	logger := logging.NewLogWithWriters(os.Stdout, os.Stderr, logging.WithVerbose())
	packClient, err := pack.NewClient(pack.WithLogger(logger))
	if err != nil {
		return fmt.Errorf("getting pack client: %v", err)
	}
	err = packClient.Build(ctx, pack.BuildOptions{
		Image:    image,
		Builder:  builder,
		AppPath:  b.source,
		Registry: "",
		Env: map[string]string{
			"GOOGLE_FUNCTION_TARGET":         b.target,
			"GOOGLE_FUNCTION_SIGNATURE_TYPE": b.funcType,
			"GOOGLE_RUNTIME":                 b.runtime,
			"X_GOOGLE_TARGET_PLATFORM":       gcfTargetPlatform,
		},
	})
	if err != nil {
		return fmt.Errorf("building function image: %v", err)
	}

	return nil
}

func (b *buildpacksFunctionServer) run() (func(), error) {
	// Create logs output files.
	var err error
	b.logStdout, err = os.Create(b.stdoutFile)
	if err != nil {
		return nil, err
	}

	b.logStderr, err = os.Create(b.stderrFile)
	if err != nil {
		return nil, err
	}
	var args = b.getDockerRunCommand()
	cmd := exec.Command(args[0], args[1:]...)
	err = cmd.Start()

	// TODO: figure out why this isn't picking up errors.
	if err != nil {
		return nil, err
	}

	// Give it some time to do its setup.
	time.Sleep(time.Duration(*startDelay) * time.Second)

	log.Printf("Framework container %q started.", b.containerID())

	return func() {
		if err := b.logs(); err != nil {
			log.Fatalf("getting container logs: %v", err)
		}
		if err := cmd.Process.Kill(); err != nil {
			log.Fatalf("failed to kill process: %v", err)
		}
		if err := b.killContainer(); err != nil {
			log.Fatalf("failed to kill container: %v", err)
		}
		log.Printf("Framework server shut down. Wrote logs to %v and %v.", b.stdoutFile, b.stderrFile)
	}, nil
}

func (b *buildpacksFunctionServer) getDockerRunCommand() []string {
	runtimeVars := []string{"docker",
		"run",
		"--network=host",
		// TODO: figure out why these aren't getting set in the buildpack.
		"--env=FUNCTION_TARGET=" + b.target,
		"--env=FUNCTION_SIGNATURE_TYPE=" + b.funcType}

	for _, s := range b.envs {
		if s != "" {
			runtimeVars = append(runtimeVars, fmt.Sprintf("--env=%s", s))
		}
	}

	return append(runtimeVars, image)
}

func (b *buildpacksFunctionServer) containerID() string {
	if b.ctID != "" {
		return b.ctID
	}
	cmd := exec.Command("docker", "ps", "--latest", "--format", "{{.ID}}")
	containerID, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to get container ID: %v", err)
	}
	b.ctID = string(bytes.TrimSpace(containerID))
	return b.ctID
}

func (b *buildpacksFunctionServer) killContainer() error {
	cmd := exec.Command("docker", "kill", b.containerID())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to kill the container %q: %v: %s", b.containerID(), err, string(output))
	}
	return nil
}

func (b *buildpacksFunctionServer) logs() error {
	defer b.logStdout.Close()
	defer b.logStderr.Close()

	args := []string{"docker", "logs", b.containerID()}
	logsCmd := exec.Command(args[0], args[1:]...)
	logsCmd.Stdout = b.logStdout
	logsCmd.Stderr = b.logStderr

	err := logsCmd.Run()
	if err != nil {
		log.Fatalf("failed to retrieve container logs: %v", err)
	}
	return nil
}
