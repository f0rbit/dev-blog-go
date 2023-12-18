#!/bin/bash
ROOT_DIR=$(pwd)
BINARY_NAME=server.out
COVERAGE_DIR=src/coverage

# Check if server.out exists
if [ ! -f "$BINARY_NAME" ]; then
    # Build the server with coverage
    echo "Building $BINARY_NAME with coverage..."
    cd src && go build -o ${BINARY_NAME} -cover && mv ${BINARY_NAME} ../
fi

cd ${ROOT_DIR}
mkdir -p ${COVERAGE_DIR}

# Start the Go server in the background
GOCOVERDIR=${COVERAGE_DIR} ./${BINARY_NAME} 2> server.log &

# Store the process ID of the Go server
server_pid=$!

# Run the tests using "bun test"
cd test && bun test

# Capture the exit code of "bun test"
test_exit_code=$?

# Stop the Go server
kill $server_pid

# Return the exit code of "bun test"
exit $test_exit_code
