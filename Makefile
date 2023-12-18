BINARY_NAME=server.out
COVERAGE_DIR=src/coverage
DATABASE_DIR=db/sqlite.db

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
	DATABASE=${DATABASE_DIR} ./${BINARY_NAME}

run-coverage: clean build-coverage
	mkdir -p ${COVERAGE_DIR}
	GOCOVERDIR=${COVERAGE_DIR} DATABASE=${DATABASE_DIR} ./${BINARY_NAME}

test: clean
	./test.sh

database:
	mkdir -p db
	touch ${DATABASE_DIR}
	sqlite3 ${DATABASE_DIR} < sql/setup.sql
	./apply_migrations.sh ${DATABASE_DIR}
	sqlite3 ${DATABASE_DIR} < sql/base_seed.sql

reset-database:
	rm -rf db
	mkdir -p db
	touch ${DATABASE_DIR}
	sqlite3 ${DATABASE_DIR} < sql/setup.sql
	./apply_migrations.sh ${DATABASE_DIR}
	sqlite3 ${DATABASE_DIR} < sql/base_seed.sql


