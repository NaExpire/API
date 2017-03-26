FROM golang:1.7.5

RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]