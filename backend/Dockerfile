# Use Go image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

COPY vendor/ vendor/
ENV GOPROXY=off
ENV GOFLAGS=-mod=vendor
ENV TZ=Asia/Shanghai

# ENV GOFLAGS=-mod=vendor

# Download dependencies
# RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
# RUN go build -o main ./cmd/main.go
RUN mkdir -p /app/upload/avatars && \
    wget -O /app/upload/avatars/default.jpg https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png

# Expose the application port
EXPOSE 8080

CMD ["sh", "-c", "go run tools/import_recipe_food_data/import_recipe_food_data.go && cd cmd && go run main.go"]