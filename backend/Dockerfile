# ---- Builder ----
FROM golang:1.24-alpine AS builder

# Install git (for go mod)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files (must exist!)
COPY go.mod go.sum ./

# Download deps
RUN go mod download

# Copy source
COPY . .

# Build (CGO=0 = safe for alpine)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ome-tv-backend .

# ---- Runtime ----
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/ome-tv-backend .

EXPOSE 8080

CMD ["./ome-tv-backend"]