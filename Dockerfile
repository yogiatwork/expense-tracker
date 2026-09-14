# ==========================================
# STAGE 1: Build the statically linked binary
# ==========================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Leverage Docker layer caching for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build flags explained:
# CGO_ENABLED=0   -> Creates a fully statically linked binary compatible with Alpine's musl libc
# -ldflags="-w -s"-> Strips debug information and symbols to keep the binary small
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/quicksnackcafe .

# ==========================================
# STAGE 2: Alpine Runtime Environment
# ==========================================
FROM alpine:latest

# 1. Install critical production system dependencies
# 2. Create a secure, locked-down non-root application user
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S quicksnackgroup && \
    adduser -S quicksnack -G quicksnackgroup

WORKDIR /home/appuser

# Copy the compiled binary from the builder stage
COPY --from=builder /app/quicksnackcafe ./quicksnackcafe

# Ensure the non-root user owns the application binary
RUN chown quicksnack:quicksnackgroup ./qucksnackcafe

# Document the port your app listens to
EXPOSE 8080

# Switch away from root to the secure user
USER quicksnack

# Exec form of ENTRYPOINT correctly passes Unix signals for graceful shutdown
ENTRYPOINT ["./quicksnackcafe"]

