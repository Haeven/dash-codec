# Start from a Go base image
FROM golang:1.20-alpine AS builder


# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

# RUN apt-get install libvpx-dev
# Update package list and install dependencies
RUN apt-get update && \
    apt-get install -y \
    build-essential \
    git \
    cmake \
    libvpx-dev \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Install vpxenc (VP9 encoder)
RUN cd /tmp && \
    git clone https://github.com/webmproject/libvpx.git && \
    cd libvpx && \
    ./configure && \
    make && \
    make install && \
    ldconfig

# Install Av1an
RUN cd /tmp && \
    git clone https://github.com/master-of-zen/Av1an.git && \
    cd Av1an && \
    make && \
    cp av1an /usr/local/bin/

# Clean up
RUN rm -rf /tmp/*

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Start a new stage from scratch
FROM alpine:latest  

# Install FFmpeg
RUN apk add --no-cache ffmpeg

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]