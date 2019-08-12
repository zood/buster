#! /bin/bash

set -e

go build -i -o buster
./buster --resources ./resources/ \
         --dev true
