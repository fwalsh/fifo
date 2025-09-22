# ---- Build stage ----
FROM golang:1.21 AS builder
WORKDIR /app

# Copy Go modules first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---- Runtime stage ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
ENTRYPOINT ["/app/server"]
