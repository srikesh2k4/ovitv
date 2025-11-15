# ---- BUILD STAGE ----
FROM golang:1.23-bullseye AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y gcc g++ sqlite3 libsqlite3-dev

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the backend source code
COPY . .

# Build Go binary with CGO enabled
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o server .

# ---- RUNTIME STAGE ----
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary
COPY --from=builder /app/server .

# Copy templates
COPY --from=builder /app/templates ./templates

# Create data folder
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./server"]
