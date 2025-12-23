# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install CA certs (needed for HTTPS calls if any)
RUN apk add --no-cache ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
# Change ./cmd/server if your main package lives elsewhere
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/server


# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Expose the port your ConnectRPC server listens on
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/server"]
