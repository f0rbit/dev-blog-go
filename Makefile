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
	rm -f server.log

run: clean build
	DATABASE=db/sqlite.db ./${BINARY_NAME}

run-coverage: clean build-coverage
	mkdir -p ${COVERAGE_DIR}
	GOCOVERDIR=${COVERAGE_DIR} DATABASE=db/sqlite.db ./${BINARY_NAME}

test: clean
	./test.sh