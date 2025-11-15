# ---------------------------
# 1. Build Stage
# ---------------------------
FROM golang:1.23-bullseye AS builder

WORKDIR /app

# Install sqlite libraries for CGO build
RUN apt-get update && apt-get install -y gcc g++ sqlite3 libsqlite3-dev

# Copy Go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of backend source
COPY . .

# Build Go binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server .

# ---------------------------
# 2. Run Stage
# ---------------------------
FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 libsqlite3-0 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates

RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./server"]
