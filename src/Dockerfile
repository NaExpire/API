FROM golang:1.7.5

RUN mkdir -p $GOPATH/src/app
ADD . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN go build -o main main.go
CMD ["/go/src/app/main"]