#! /bin/bash

set -e

go build -o buster
./buster --resources ./assets/ \
         --dev true
