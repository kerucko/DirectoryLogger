FROM golang:1.20-alpine

WORKDIR /usr/src/directory_logger

RUN apk --no-cache add bash make

COPY go.mod go.sum ./
RUN go mod download \
    && go mod verify \
    && go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ./ ./
RUN make build