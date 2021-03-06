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
	"os/exec"
	"syscall"
)

func newCmd(args []string) *exec.Cmd {
	cmd := exec.Command(args[0], args[1:]...)

	// Make a process group so that later we can kill child processes too. As an
	// example, if the command is `go run main.go`, Go will build a binary in a
	// temp dir and then execute it. If we simply cmd.Process.Kill() the exec.Command
	// then the running binary will not be killed. Only if we make a group will
	// child processes be killed.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	return cmd
}

func stopCmd(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}
	return nil
}
