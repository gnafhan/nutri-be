#!/bin/bash

# Colors for terminal output
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color

# Function to display messages
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if .env file exists, if not create from example
if [ ! -f .env ]; then
    info "Creating .env file from .env.example"
    cp .env.example .env
    warning "Please check and update the .env file with your configuration"
    sleep 2
fi

# Function to start containers
start_containers() {
    info "Starting Nutribox API containers..."
    docker-compose up -d
    if [ $? -eq 0 ]; then
        info "Nutribox API is now running!"
        info "API is accessible at: http://localhost:$(grep APP_PORT .env | cut -d '=' -f2 || echo 3000)"
    else
        error "Failed to start containers. Please check the logs."
        exit 1
    fi
}

# Function to stop containers
stop_containers() {
    info "Stopping Nutribox API containers..."
    docker-compose down
    if [ $? -eq 0 ]; then
        info "Containers stopped successfully."
    else
        error "Failed to stop containers."
        exit 1
    fi
}

# Function to rebuild and restart containers
rebuild_containers() {
    info "Rebuilding and restarting Nutribox API containers..."
    docker-compose down
    docker-compose up -d --build
    if [ $? -eq 0 ]; then
        info "Containers rebuilt and restarted successfully."
        info "API is accessible at: http://localhost:$(grep APP_PORT .env | cut -d '=' -f2 || echo 3000)"
    else
        error "Failed to rebuild containers."
        exit 1
    fi
}

# Function to show logs
show_logs() {
    info "Showing logs for Nutribox API containers..."
    docker-compose logs -f
}

# Function to show help
show_help() {
    echo -e "\nNutribox API Docker Management Script\n"
    echo -e "Usage: $0 [command]\n"
    echo -e "Commands:"
    echo -e "  start\t\tStart the containers"
    echo -e "  stop\t\tStop the containers"
    echo -e "  restart\tRestart the containers"
    echo -e "  rebuild\tRebuild and restart the containers"
    echo -e "  logs\t\tShow container logs"
    echo -e "  help\t\tShow this help message\n"
    echo -e "If no command is provided, containers will be started by default.\n"
}

# Make script executable
chmod +x "$0"

# Process command line arguments
case "$1" in
    start)
        start_containers
        ;;
    stop)
        stop_containers
        ;;
    restart)
        stop_containers
        start_containers
        ;;
    rebuild)
        rebuild_containers
        ;;
    logs)
        show_logs
        ;;
    help)
        show_help
        ;;
    *)
        # Default action if no argument is provided
        if [ -z "$1" ]; then
            start_containers
        else
            error "Unknown command: $1"
            show_help
            exit 1
        fi
        ;;
esac

exit 0