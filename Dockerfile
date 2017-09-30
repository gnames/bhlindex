FROM golang:1.9

ENV LAST_FULL_REBUILD 2017-10-01

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get -u -d github.com/mattes/migrate/cli github.com/lib/pq github.com/satori/go.uuid
RUN go get -u -d github.com/GlobalNamesArchitecture/gnfinder
# RUN cd $GOPATH/src/github.com/GlobalNamesArchitecture/gnfinder && git pull && cd -
RUN go build -tags 'postgres' -o /go/bin/migrate github.com/mattes/migrate/cli

RUN apt-get update && apt-get -yq install postgresql-client

WORKDIR /go/src/github.com/GlobalNamesArchitecture/bhlindex
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

ENTRYPOINT scripts/development.sh
