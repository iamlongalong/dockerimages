#!/bin/bash

set -e
set -x

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    rm -rf test_output*
    rm -f dockerimages
    docker rmi -f redis:alpine busybox:latest || true
}

# Setup
echo "Setting up test environment..."
mkdir -p test_output{,2,3,4}
trap cleanup EXIT

# Build the binary
echo "Building dockerimages binary..."
(cd .. && go build -o example/dockerimages)
if [ ! -f "./dockerimages" ]; then
    echo "Failed to build binary"
    exit 1
fi

# Test 1: Basic usage with multiple images
echo "Test 1: Basic usage with multiple images..."
./dockerimages save -o test_output busybox:latest redis:alpine
[ -f "test_output/busybox_latest.tar" ] && echo "✅ busybox image saved successfully" || exit 1
[ -f "test_output/redis_alpine.tar" ] && echo "✅ redis image saved successfully" || exit 1

# Test 2: Using image list file
echo "Test 2: Using image list file..."
./dockerimages save -o test_output2 -f images.txt
[ -f "test_output2/busybox_latest.tar" ] && echo "✅ busybox image from file saved successfully" || exit 1
[ -f "test_output2/redis_alpine.tar" ] && echo "✅ redis image from file saved successfully" || exit 1

# Test 3: Using docker-compose file
echo "Test 3: Using docker-compose file..."
./dockerimages save -o test_output3 -c docker-compose.yml
[ -f "test_output3/busybox_latest.tar" ] && echo "✅ busybox image from compose file successfully" || exit 1
[ -f "test_output3/redis_alpine.tar" ] && echo "✅ redis image from compose file successfully" || exit 1

# Test 4: Merge images
echo "Test 4: Testing merge functionality..."
./dockerimages save -o test_output4 -m busybox:latest redis:alpine
[ -f "test_output4/images.tar" ] && echo "✅ merged images saved successfully" || exit 1

echo "All tests completed successfully! 🎉" 