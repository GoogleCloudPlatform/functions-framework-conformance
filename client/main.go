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

// This binary contains a valiation framework for functions frameworks.
package main

import (
	"flag"
	"log"
	"strings"
)

var (
	runCmd = flag.String("cmd", "", "string with command to run a Functions Framework server at localhost:8080. Ignored if -buildpacks=true.")
	// functionSignature is the function's signature as signature in GCF i.e. will be set in the `GOOGLE_FUNCTION_SIGNATURE_TYPE` env variable.
	functionSignature = flag.String("type", "http", "the function signature to use (must be 'http', 'cloudevent', or 'legacyevent'")
	// declarativeSignature indicates the declarative function signature that is being tested. This is used to test `typed` functions which are exposed to GCF as the `http` signature type.
	declarativeSignature    = flag.String("declarative-type", "", "the declarative signature type of the function (must be 'http', 'cloudevent', 'legacyevent', or 'typed'), default matches -type")
	validateMapping         = flag.Bool("validate-mapping", true, "whether to validate mapping from legacy->cloud events and vice versa (as applicable)")
	outputFile              = flag.String("output-file", "function_output.json", "name of file output by function")
	useBuildpacks           = flag.Bool("buildpacks", true, "whether to use the current release of buildpacks to run the validation. If true, -cmd is ignored and --builder-* flags must be set.")
	source                  = flag.String("builder-source", "", "function source directory to use in building. Required if -buildpacks=true")
	target                  = flag.String("builder-target", "", "function target to use in building. Required if -buildpacks=true")
	runtime                 = flag.String("builder-runtime", "", "runtime to use in building. Required if -buildpacks=true")
	tag                     = flag.String("builder-tag", "latest", "builder image tag to use in building")
	startDelay              = flag.Uint("start-delay", 1, "Seconds to wait before sending HTTP request to command process")
	validateConcurrencyFlag = flag.Bool("validate-concurrency", false, "whether to validate concurrent requests can be handled, requires a function that sleeps for 1 second ")
	envs                    = flag.String("envs", "", "a comma separated string of additional runtime environment variables")
)

func main() {
	flag.Parse()

	if *useBuildpacks {
		if *runtime == "" || *source == "" || *target == "" {
			log.Fatalf("testing via buildpacks requires -builder-runtime, -builder-source, and -builder-target to be set")
		}
	}

	if *declarativeSignature == "" {
		*declarativeSignature = *functionSignature
	}

	// Set runtime env vars that reflect https://cloud.google.com/functions/docs/configuring/env-var
	validationRuntimeEnv := []string{"FUNCTION_SIGNATURE_TYPE=" + *functionSignature}
	validationRuntimeEnv = append(validationRuntimeEnv, strings.Split(*envs, ",")...)

	v := newValidator(validatorParams{
		validateMapping:      *validateMapping,
		useBuildpacks:        *useBuildpacks,
		runCmd:               *runCmd,
		outputFile:           *outputFile,
		source:               *source,
		target:               *target,
		runtime:              *runtime,
		functionSignature:    *functionSignature,
		declarativeSignature: *declarativeSignature,
		tag:                  *tag,
		validateConcurrency:  *validateConcurrencyFlag,
		envs:                 validationRuntimeEnv,
	})

	if err := v.runValidation(); err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("All validation passed!")
}
