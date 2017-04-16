FROM golang:1.8.1

RUN mkdir -p /app
RUN mkdir -p $GOPATH/src/github.com/NAExpire/API/src
ADD src $GOPATH/src/github.com/NAExpire/API/src
WORKDIR $GOPATH/src/github.com/NAExpire/API/src
RUN go build -o /app/API
EXPOSE 8000
ENTRYPOINT ["/app/API"]
