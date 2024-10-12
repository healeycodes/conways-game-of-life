#!/bin/bash

# Install X11 development libraries
sudo apt-get update -y
sudo apt-get install -y libx11-dev libxrandr-dev libxcursor-dev libxi-dev mesa-common-dev libgl1-mesa-dev libglu1-mesa-dev

# Build the Go application
go build -tags netgo -ldflags '-s -w' -o app .