# Doesn't work, just keep as a relic

FROM ubuntu:latest

# Set working directory
WORKDIR /app

# Install Go and update CA certificates
RUN apt-get update -y && \
  apt-get install -y golang-go ca-certificates

# Copy Go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install X11 development libraries (for Ebiten)
RUN apt-get update -y && \
  apt-get install -y libx11-dev libxrandr-dev libxcursor-dev libxi-dev libxinerama-dev mesa-common-dev libgl1-mesa-dev libglu1-mesa-dev

# Copy application code
COPY . .

# Build the application
RUN go build -tags netgo -ldflags '-s -w' -o app .

# Expose port (if necessary)
EXPOSE 8080

# Run the command to start the application
CMD ["./app"]