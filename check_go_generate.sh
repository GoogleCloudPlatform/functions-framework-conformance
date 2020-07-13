#!/bin/bash

set -e

go generate ./...
if ! git diff --name-status --exit-code HEAD;
then
    echo ERROR: Please run '"go generate ./..."' after making changes to files under events/generate/data.
    exit 1
fi;
