ARG work_dir=/go/src/github.com/Bnei-Baruch/mdb
ARG build_number=dev
ARG mdb_url="postgres://user:password@host.docker.internal:5435/mdb?sslmode=disable"
ARG mdb_test_url="postgres://user:password@host.docker.internal:5435/%s?sslmode=disable"

FROM golang:1.17-alpine3.15 as build

LABEL maintainer="edoshor@gmail.com"

ARG work_dir
ARG build_number
ARG mdb_url
ARG mdb_test_url

ENV GOOS=linux \
	CGO_ENABLED=0 \
	MDB_URL=${mdb_url} \
	MDB_TEST_URL=${mdb_test_url}

RUN apk update && \
    apk add --no-cache \
    git

WORKDIR ${work_dir}
COPY . .
RUN go test -v $(go list ./... | grep -v /models) \
    && go build -ldflags "-w -X github.com/Bnei-Baruch/mdb/version.PreRelease=${build_number}"

FROM alpine:3.15

ARG work_dir
WORKDIR /app
COPY misc/*.sh ./
COPY --from=build ${work_dir}/mdb .

EXPOSE 8080
CMD ["./mdb", "server"]
