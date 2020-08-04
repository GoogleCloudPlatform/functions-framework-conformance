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
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/buildpacks/pack"
)

const (
	image = "conformance-test-func"
)

type buildpacksFunctionServer struct {
	output      string
	source      string
	target      string
	funcType    string
	runtime     string
	tag         string
	containerID string
}

func (b *buildpacksFunctionServer) Start() (func(), error) {
	ctx := context.Background()
	typ := *functionType
	if typ == "legacyevent" {
		typ = "event"
	}

	if err := b.build(ctx); err != nil {
		return nil, fmt.Errorf("building function container: %v", err)
	}

	shutdown, err := b.run(ctx)
	if err != nil {
		return nil, fmt.Errorf("running function container: %v", err)
	}

	return shutdown, nil
}

func (b *buildpacksFunctionServer) OutputFile() ([]byte, error) {
	contents, _, err := b.dockerClient.CopyFromContainer(context.Background(), b.containerID, "/workspace/"+b.output)
	if err != nil {
		return nil, fmt.Errorf("fetching function output from container: %v", err)
	}
	// Write to local dir for debugging.
	f, err := os.Create(b.output)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("creating function output file locally: %v", err)
	}
	if _, err := io.Copy(f, contents); err != nil {
		return nil, fmt.Errorf("writing function output file locally: %v", err)
	}
	return ioutil.ReadAll(contents)
}

func (b *buildpacksFunctionServer) build(ctx context.Context) error {
	// TODO: use latest tag once GCF builders have it
	builder := fmt.Sprintf("us.gcr.io/fn-img/buildpacks/%s/builder:%s", b.runtime, b.tag)
	packClient, err := pack.NewClient()
	if err != nil {
		return fmt.Errorf("getting pack client: %v", err)
	}
	err = packClient.Build(ctx, pack.BuildOptions{
		Image:    image,
		Builder:  builder,
		AppPath:  b.source,
		Registry: "",
		Env: map[string]string{
			"GOOGLE_FUNCTION_SOURCE":         b.source,
			"GOOGLE_FUNCTION_TARGET":         b.target,
			"GOOGLE_FUNCTION_SIGNATURE_TYPE": b.funcType,
		},
	})
	if err != nil {
		return fmt.Errorf("building function image: %v", err)
	}

	return nil
}

func (b *buildpacksFunctionServer) run(ctx context.Context) (func(), error) {

	return func() {

		log.Printf("Framework server shut down. Wrote logs to %v and %v.", stdoutFile, stderrFile)
	}, nil
}
