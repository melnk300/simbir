FROM golang:latest AS build
WORKDIR /go/src/rest
COPY . .
RUN go build -o /simbir ./cmd/simbir/main.go