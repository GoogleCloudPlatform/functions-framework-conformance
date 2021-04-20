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
	"testing"
)

func TestStartAndShutdown(t *testing.T) {
	server := localFunctionServer{
		// Use a command that execs another command in order to test that the whole
		// process group is killed.
		cmd: "/bin/sh -c 'exec sleep 90'",
	}

	shutdown, err := server.Start()
	if shutdown != nil {
		defer shutdown()
	}

	if err != nil {
		t.Errorf("unable to start localFunctionServer: %v", err)
	}
}
