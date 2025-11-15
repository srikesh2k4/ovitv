# ---------------------
#       BUILDER
# ---------------------
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Enable CGO (sqlite3 requires it)
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Build binary
RUN go build -o server .

# ---------------------
#       RUNTIME
# ---------------------
FROM alpine:latest
RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
