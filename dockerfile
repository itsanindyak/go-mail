# syntax=docker/dockerfile:1.6

############################################
# 1️⃣ Builder Stage
############################################
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine AS builder

# Install required packages
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Enable Go Modules
ENV GO111MODULE=on

# Copy go mod files first (better caching)
COPY go.mod go.sum ./

# Download dependencies (cached unless mod files change)
RUN go mod download

# Copy rest of the source code
COPY . .

# Build arguments for cross compilation
ARG TARGETOS
ARG TARGETARCH

# Build optimized static binary
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o /app/app ./cmd


############################################
# 2️⃣ Minimal Runtime Stage
############################################
FROM gcr.io/distroless/base-debian12:nonroot

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/app .

# Expose port (change if needed)
EXPOSE 8080

# Run as non-root (distroless already uses nonroot user)
USER nonroot:nonroot

ENTRYPOINT ["/app/app"]
