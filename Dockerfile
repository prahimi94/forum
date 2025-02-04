# Use the official Go 1.23 image to build the application
FROM golang:1.23-alpine AS builder


LABEL version="1.0" \
        description="Forum Dockerized Application"\
        author="Parisa Rahimi \ Mahdi Kheirkhah \ Fatemeh Kheirkhah \ Majid Rouhani" \ 
        environment="development"

# Enable CGO and install required libraries
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Install dependencies
RUN apk add --no-cache gcc musl-dev sqlite sqlite-libs

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -tags "sqlite3" -o forum .

# Use Alpine as the final image
FROM alpine:latest

# Set the working directory in the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/forum .
COPY --from=builder /app/db ./db
COPY --from=builder /app/modules ./modules
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/utils ./utils

# Install SQLite and required libraries for the database
RUN apk add --no-cache sqlite sqlite-libs

# Expose the port your application runs on
EXPOSE 8080

# Command to run the executable
CMD ["./forum"]
