GO_FILES      = $(shell find . -path ./vendor -prune -o -type f -name "*.go" -print)
IMPORT_PATH   = $(shell pwd | sed "s|^$(GOPATH)/src/||g")
GIT_HASH      = $(shell git rev-parse HEAD)
LDFLAGS       = -w -X $(IMPORT_PATH)/version.PreRelease=$(PRE_RELEASE)

build: clean test
	@go build -ldflags '$(LDFLAGS)'

clean:
	@rm -f mdb

install:
	@godep restore

test:
	@go test

lint:
	@golint $(GO_FILES) || true


