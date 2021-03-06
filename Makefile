GO := go
PREFIX ?= $(HOME)/.local

SERVER_BIN := myroomies-server
CLIENT_BIN := myroomies-client

.PHONY: all
all: build-server build-client

.PHONY: build-server
build-server:
	$(GO) build ./cmd/$(SERVER_BIN)

.PHONY: build-client
build-client:
	$(GO) build ./cmd/$(CLIENT_BIN)

.PHONY: install-client
install-client: build-client
	cp $(CLIENT_BIN) $(PREFIX)/bin/myroomies

.PHONY: test-e2e
test-e2e: all
	echo "Running the end-to-end tests suite" && \
	./e2e/libs/bats-core/bin/bats e2e 

.PHONY: test-e2e-mongodb
test-e2e-mongodb: all
	echo "Starting MongoDB..." && \
	docker run -d -p 27017:27017 --name myroomies-mongo mongo | true && \
	sleep 3 && \
	echo "MongoDB started. Listening on port 27017" && \
	echo "Running the end-to-end tests suite" && \
	MYROOMIES_E2E_TESTS_MONGODB_ADDRESS="mongodb://localhost:27017" ./e2e/libs/bats-core/bin/bats e2e && \
	docker rm -f myroomies-mongo

.PHONY: test-unittest
test-unittest:
	$(GO) test ./...

.PHONY: test
test: test-unittest test-e2e
