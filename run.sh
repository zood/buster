#! /bin/bash

set -e

go build -i -o buster
./buster --resources ~/coding/gocode/src/zood.xyz/buster/resources/ \
         --dev true