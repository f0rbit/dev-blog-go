BINARY_NAME=server.out
COVERAGE_DIR=src/coverage
DATABASE_DIR=db/sqlite.db
TEST_PORT=8080
TEST_TOKEN=TEST_TOKEN_123

all: clean build run

build:
	@cd src && go build -o ${BINARY_NAME} && mv ${BINARY_NAME} ../
	# code sign it (this is for macos) (would be fixed with docker?)
	# @codesign --sign - ./${BINARY_NAME}
	# @codesign --verify --verbose ./${BINARY_NAME}
	@echo "Built binary"

build-coverage:
	@cd src && go build -o ${BINARY_NAME} -cover && mv ${BINARY_NAME} ../
	# code sign it (this is for macos) (would be fixed with docker?)
	# @codesign --sign - ./${BINARY_NAME}
	# @codesign --verify --verbose ./${BINARY_NAME}
	@echo "Built binary with coverage"

clean:
	@rm -f src/${BINARY_NAME}
	@rm -f ${BINARY_NAME}
	@rm -f logs/app.log
	@rm -rf ${COVERAGE_DIR}
	@rm -f server.log
	@rm -rf client/dist
	@echo "Cleaned files"

run: clean build
	@echo "Running..."
	@DATABASE=${DATABASE_DIR}?parseTime=true ./${BINARY_NAME}

run-coverage: clean build-coverage
	@echo "Runnnig with coverage"
	@mkdir -p ${COVERAGE_DIR}
	@GOCOVERDIR=${COVERAGE_DIR} DATABASE=${DATABASE_DIR}?parseTime=true ./${BINARY_NAME}

test: clean
	@PORT=${TEST_PORT} AUTH_TOKEN=${TEST_TOKEN} ./test.sh

database:
	@mkdir -p db
	@touch ${DATABASE_DIR}
	@sqlite3 ${DATABASE_DIR} < sql/setup.sql
	@echo "Setup SQL"
	@./sql/apply_migrations.sh ${DATABASE_DIR}
	@sqlite3 ${DATABASE_DIR} < sql/base_seed.sql
	@echo "Executed base_seed.sql"

reset-database:
	rm -rf db
	mkdir -p db
	touch ${DATABASE_DIR}
	sqlite3 ${DATABASE_DIR} < sql/setup.sql
	./sql/apply_migrations.sh ${DATABASE_DIR}
	sqlite3 ${DATABASE_DIR} < sql/base_seed.sql

coverage: 
	@PORT=${TEST_PORT} AUTH_TOKEN=${TEST_TOKEN} ./test.sh > /dev/null 2> /dev/null
	@go tool covdata func -i=${COVERAGE_DIR} | grep total | awk '{print $$3}'

coverage-report:
	@PORT=${TEST_PORT} AUTH_TOKEN=${TEST_TOKEN} ./test.sh > /dev/null 2> /dev/null
	@go tool covdata func -i=${COVERAGE_DIR}

build-client:
	cd client && npm i && npm run build

run-client:
	cd client && npm run dev
