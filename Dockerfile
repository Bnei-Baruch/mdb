ARG work_dir=/go/src/github.com/Bnei-Baruch/mdb
ARG build_number=dev

FROM golang:1.24-alpine AS build

LABEL maintainer="edoshor@gmail.com"

ARG work_dir
ARG build_number

ENV GOOS=linux \
	CGO_ENABLED=0

RUN apk update && \
    apk add --no-cache \
    git

WORKDIR ${work_dir}
COPY . .
RUN go build -ldflags "-w -X github.com/Bnei-Baruch/mdb/version.PreRelease=${build_number}"

FROM alpine

RUN apk update && \
    apk add --no-cache \
    mailx \
    postfix

RUN echo "mydomain = bbdomain.org" >> /etc/postfix/main.cf
RUN echo "myhostname = app.mdb" >> /etc/postfix/main.cf
RUN echo "myorigin = \$mydomain" >> /etc/postfix/main.cf
RUN echo "relayhost = [smtp.local]" >> /etc/postfix/main.cf

ARG work_dir
WORKDIR /app
COPY misc/*.sh ./
COPY --from=build ${work_dir}/mdb .
COPY --from=build ${work_dir}/migrations migrations

EXPOSE 8080
CMD ["./mdb", "server"]
