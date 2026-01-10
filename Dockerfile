# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies for CGO (SQLite)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO enabled and FTS5 support
RUN CGO_ENABLED=1 go build -tags "fts5" -o bin/glossary ./cmd/glossary

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies for SQLite
RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/glossary .

# Create data directory
RUN mkdir -p /app/data

# Environment variables
ENV PORT=8080
ENV DATABASE_PATH=/app/data/glossary.db

EXPOSE 8080

CMD ["./glossary"]
