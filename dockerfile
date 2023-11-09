FROM golang:latest

WORKDIR /app

COPY go.mod go.mod
COPY main.go main.go
COPY commands commands

RUN go get
RUN go build -o ./bin

ENTRYPOINT [ "/app/bin" ]