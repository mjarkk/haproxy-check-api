#!/bin/bash

# remove old files
rm -rf releases

# Bulid the go project
GOOS=linux go build

# make releases folder and copy over the needed files
mkdir releases
cp haproxy-check-api releases
cp Dockerfile releases

# zip the files
cd releases
zip release.zip haproxy-check-api Dockerfile
cd ..

# Cleanup some files
rm releases/Dockerfile releases/haproxy-check-api

echo "----------------------------"
echo "Created releases/release.zip"
echo "----------------------------"
