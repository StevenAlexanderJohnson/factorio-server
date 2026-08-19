# Stage 1: Build Go binary
FROM golang:1.26-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/factorio-api ./cmd

# Stage 2: Runtime image
# Using debian:bookworm-slim because the official Factorio headless binary requires glibc and xz-utils
FROM debian:bookworm-slim

# Install dependencies required by Factorio and archive extraction
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tar \
    xz-utils \
    libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Create necessary directories for Factorio saves, data, and binary
RUN mkdir -p /opt /factorio/data/saves /app

# Copy the compiled application and default configuration
COPY --from=builder /bin/factorio-api /app/factorio-api
COPY config.yaml /app/config.yaml

# Expose Factorio default game port (UDP) and HTTP API port
EXPOSE 34197/udp
EXPOSE 8080

ENTRYPOINT ["/app/factorio-api"]
CMD ["-config", "/app/config.yaml"]
