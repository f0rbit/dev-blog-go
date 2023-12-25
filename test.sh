#!/bin/bash
ROOT_DIR=$(pwd)
BINARY_NAME=server.out
COVERAGE_DIR=src/coverage
DATABASE_FILE=db/test.db

mkdir -p db
rm -f ${DATABASE_FILE}
touch ${DATABASE_FILE}

# setup the database
sqlite3 ${DATABASE_FILE} < sql/setup.sql
sqlite3 ${DATABASE_FILE} < sql/test_seed.sql
${ROOT_DIR}/apply_migrations.sh ${DATABASE_FILE}

# Check if server.out exists
if [ ! -f "$BINARY_NAME" ]; then
    # Build the server with coverage
    echo "Building $BINARY_NAME with coverage..."
    cd src && go build -o ${BINARY_NAME} -cover && mv ${BINARY_NAME} ../
fi

cd ${ROOT_DIR}
mkdir -p ${COVERAGE_DIR}

# Start the Go server in the background
GOCOVERDIR=${COVERAGE_DIR} DATABASE=${DATABASE_FILE} ./${BINARY_NAME} 2> server.log &

# Store the process ID of the Go server
server_pid=$!

# install dependencies
cd ${ROOT_DIR}/client && ~/.bun/bin/bun i
cd ${ROOT_DIR}/test && ~/.bun/bin/bun i
cd ${ROOT_DIR}

# Run the tests using "bun test"
cd test && ~/.bun/bin/bun test

# Capture the exit code of "bun test"
test_exit_code=$?

# Stop the Go server
kill $server_pid

# Return the exit code of "bun test"
exit $test_exit_code
