GO_FILES      = $(shell find . -path ./vendor -prune -o -type f -name "*.go" -print)
GIT_HASH      = $(shell git rev-parse HEAD)
LDFLAGS       = -w -X main.commitHash=$(GIT_HASH)

build: clean
	@go build -ldflags '$(LDFLAGS)'

clean:
	@rm -f mdb

install:
	@godep restore

test:
	@go test

lint:
	@golint $(GO_FILES) || true

try:
	@echo $(GO_FILES)

