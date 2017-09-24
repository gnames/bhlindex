FROM golang:1.9

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get -u -d github.com/mattes/migrate/cli github.com/lib/pq github.com/satori/go.uuid
RUN go build -tags 'postgres' -o /go/bin/migrate github.com/mattes/migrate/cli

RUN apt-get update && apt-get -yq install postgresql-client

WORKDIR /go/src/github.com/GlobalNamesArchitecture/bhlindex
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

VOLUME "/tmp/gni_mysql"

ENTRYPOINT scripts/development.sh
