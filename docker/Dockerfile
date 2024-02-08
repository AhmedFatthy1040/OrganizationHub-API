# Use the official GoLang base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod .
COPY go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the entire project directory into the container
COPY . .

# Build the GoLang application
RUN go build -o main .

# Command to run the application
CMD ["./main"]