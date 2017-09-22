FROM golang:1.9

WORKDIR /go/src/github.com/GlobalNamesArchitecture/bhlindex
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega

VOLUME "/tmp/gni_mysql"

ENTRYPOINT ["ginkgo", "watch"]
