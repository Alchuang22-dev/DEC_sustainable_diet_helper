# Use Go image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

RUN apt-get install -y ca-certificates && update-ca-certificates

# Copy the Go module files
COPY go.mod go.sum ./
COPY .env .env

COPY vendor/ vendor/
ENV GOPROXY=off
ENV GOFLAGS=-mod=vendor

# ENV GOFLAGS=-mod=vendor

# Download dependencies
# RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/main.go

# Expose the application port
EXPOSE 8080

CMD ["sh", "-c", "go run tools/import_recipe_food_data/import_recipe_food_data.go && ./main"]