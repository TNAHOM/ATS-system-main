# Use official Go image to build the binary
FROM golang:1.24.5 as builder

# Set work directory inside builder
WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code into container
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -o aiServer ./cmd/main.go


# Use minimal final image for security and performance
FROM alpine:3.20

# Set working directory inside final image
WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/aiServer .

# Expose gateway port
EXPOSE 8081

# Run the gateway binary
CMD ["/app/aiServer"]
