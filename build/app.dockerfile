# ========== STAGE 1: BUILD ==========
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy dependency files terlebih dahulu
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .


# Build binary untuk worker
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# ========== STAGE 2: RUN ==========
FROM alpine:latest

# Install dependency minimal (optional)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy hasil build dari stage builder
COPY --from=builder /app/app .
COPY --from=builder /app/.env .

# Expose port jika worker perlu komunikasi HTTP (optional)
# EXPOSE 8080

# Jalankan worker
CMD ["./app"]
