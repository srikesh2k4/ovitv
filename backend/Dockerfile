# ---- Build Go Binary ----
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first (for Docker caching)
COPY go.mod go.sum ./

# Download modules
RUN go mod download

# Copy the entire backend source
COPY . .

# Build binary (production safe)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ome-tv-backend .

# ---- Runtime ----
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy backend binary
COPY --from=builder /app/ome-tv-backend .

EXPOSE 8080

CMD ["./ome-tv-backend"]
