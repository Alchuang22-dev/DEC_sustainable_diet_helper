# Use Go image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/main.go

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]