FROM golang:1.9

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."
VOLUME "/tmp/gni_mysql"

ENTRYPOINT ["go-wrapper", "run"]
