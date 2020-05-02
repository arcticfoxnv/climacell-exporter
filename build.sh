#!/bin/bash

sed -i -e "s/Commit.*/Commit = \"$(git rev-parse --short=8 HEAD)\"/" version.go

go install .
