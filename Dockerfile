# Use the official Go image as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod download github.com/bytedance/sonic

# Copy the entire source code to the working directory
COPY . .

# Build the Go app
RUN go build main.go

# Expose the port the app runs on
EXPOSE 8080

# Run the compiled binary
CMD ["./main"]
