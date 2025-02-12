.PHONY: test-setup test test-cleanup build run

build:
	sudo docker-compose build

run:
	sudo docker-compose up

test-setup:
	sudo ./scripts/setup_test_db.sh

test: test-setup
	TEST_DB_HOST=localhost \
	TEST_DB_PORT=5433 \
	TEST_DB_USER=postgres \
	TEST_DB_PASSWORD=postgres123 \
	TEST_DB_NAME=neoway_test \
	go test ./... -v

test-cleanup:
	sudo docker-compose -f docker-compose.test.yml down -v

test-all: test-cleanup test-setup test test-cleanup

down:
	sudo docker-compose down -v 