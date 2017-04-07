FROM golang:1.7.5

RUN mkdir -p /app
RUN mkdir -p $GOPATH/src/github.com/NAExpire/API/src
ADD . $GOPATH/src/github.com/NAExpire/API/src
WORKDIR $GOPATH/src/github.com/NAExpire/API/src
RUN go build -o /app/API main.go
EXPOSE 8000
ENTRYPOINT ["/app/API"]
