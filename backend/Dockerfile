# ---- BUILD STAGE ----
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend source
COPY . .

# Build the Go binary
RUN go build -o server .

# ---- RUN STAGE ----
FROM alpine:latest

RUN apk add --no-cache sqlite-libs

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .

# Copy templates folder
COPY --from=builder /app/templates ./templates

# Create data folder inside container
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./server"]
