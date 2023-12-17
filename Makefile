BINARY_NAME=server.out
COVERAGE_DIR=src/coverage

all: clean build run

build: 
	cd src && go build -o ${BINARY_NAME} && mv ${BINARY_NAME} ../

build-coverage:
	cd src && go build -o ${BINARY_NAME} -cover && mv ${BINARY_NAME} ../

clean:
	rm -f src/${BINARY_NAME}
	rm -f ${BINARY_NAME}
	rm -f logs/app.log
	rm -rf ${COVERAGE_DIR}

run: clean build
	./${BINARY_NAME}

run-coverage: clean build-coverage
	GOCOVERDIR=${COVERAGE_DIR} ./${BINARY_NAME}

test: build-coverage
	cd test && bun test
