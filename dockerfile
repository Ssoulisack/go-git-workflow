# Build Stage
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git

WORKDIR /go/src/app

# Copy go mod files first and download deps
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the application
COPY . .

# Build the application
RUN go build -o /go/bin/app -v ./cmd/main.go