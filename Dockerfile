# ใช้ base image
FROM golang:1.24.3-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build binary from cmd/main.go
RUN go build -o main ./cmd

# Expose port (เช่นถ้าใช้ 8080)
EXPOSE 8080

# Run the application
CMD ["./main"]
