FROM golang:1.11

ENV LAST_FULL_REBUILD 2018-11-17

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get -u github.com/hashicorp/go-multierror
RUN go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq
RUN go build -tags 'postgres' -o /go/bin/migrate github.com/golang-migrate/migrate/cli
#
RUN apt-get update && apt-get -yq install postgresql-client

WORKDIR /bhlindex
COPY . .
ENV GO111MODULE on

ENTRYPOINT bhlindex/development.sh
