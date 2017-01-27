GO_FILES      = $(shell find . -path ./vendor -prune -o -type f -name "*.go" -print)
IMPORT_PATH   = $(shell pwd | sed "s|^$(GOPATH)/src/||g")
GIT_HASH      = $(shell git rev-parse HEAD)
LDFLAGS       = -w -X $(IMPORT_PATH)/version.PreRelease=$(PRE_RELEASE)
APIB_FILES    = $(shell find . -type f -path "./*/*.apib")

build: clean test api
	@go build -ldflags '$(LDFLAGS)'

clean:
	rm -f mdb

install:
	@godep restore

test:
	go test $(shell go list ./... | grep -v /vendor/)

lint:
	@golint $(GO_FILES) || true

fmt:
	@gofmt -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

api: test
	cp docs.tmpl docs.apib; \
	for f in ${APIB_FILES}; \
	do \
	cat $$f >> docs.apib; \
	done; \
	aglio -i docs.apib -o docs.html --theme-template triple
