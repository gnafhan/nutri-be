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

# Function to check Docker is installed and running
check_docker() {
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! docker info &> /dev/null; then
        error "Docker daemon is not running. Please start Docker first."
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi

    info "Docker and Docker Compose are properly installed and running."
}

# Check if .env file exists, if not create from example
if [ ! -f .env ]; then
    info "Creating .env file from .env.example"
    cp .env.example .env
    warning "Please check and update the .env file with your configuration"
    sleep 2
fi

# Function to check container status
check_container_status() {
    container_name=$1
    if [ "$(docker ps -q -f name=$container_name)" ]; then
        info "$container_name is running."
        return 0
    elif [ "$(docker ps -aq -f status=exited -f name=$container_name)" ]; then
        warning "$container_name exists but is not running."
        info "Checking container logs..."
        docker logs $container_name | tail -n 20
        return 1
    else
        warning "$container_name does not exist."
        return 2
    fi
}

# Function to start containers
start_containers() {
    info "Starting Nutribox API containers..."
    docker-compose up -d
    if [ $? -eq 0 ]; then
        info "Containers started. Checking status..."
        sleep 5
        
        # Check each container
        check_container_status "nutribox-api"
        api_status=$?
        check_container_status "nutribox-postgres"
        db_status=$?
        
        if [ $api_status -eq 0 ] && [ $db_status -eq 0 ]; then
            info "Nutribox API is now running!"
            info "API is accessible at: http://localhost:$(grep APP_PORT .env | cut -d '=' -f2 || echo 3000)"
        else
            warning "Some containers may not be running properly."
            info "You can check detailed logs with: './run-docker.sh logs'"
        fi
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

# Function to troubleshoot common issues
troubleshoot() {
    info "Running diagnostics on Nutribox API containers..."
    
    # Check if containers are running
    info "Checking container status..."
    docker ps -a | grep nutribox
    
    # Check network
    info "Checking Docker networks..."
    docker network ls | grep nutribox
    
    # Check environment variables
    info "Checking environment variables..."
    docker-compose config | grep -E "APP_PORT|DB_HOST|DB_PORT"
    
    # Check port bindings
    info "Checking port bindings..."
    docker-compose ps
    
    # Suggestion
    info "Diagnostic complete. If you're still experiencing issues:"
    info "1. Try 'docker-compose down -v' to remove volumes and restart"
    info "2. Check if the ports are already in use on your system"
    info "3. Verify the database settings in .env file match docker-compose.yml"
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
    echo -e "  status\tCheck the status of containers"
    echo -e "  troubleshoot\tRun diagnostics on containers"
    echo -e "  help\t\tShow this help message\n"
    echo -e "If no command is provided, containers will be started by default.\n"
}

# Make script executable
chmod +x "$0"

# Process command line arguments
case "$1" in
    start)
        check_docker
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
        check_docker
        rebuild_containers
        ;;
    logs)
        show_logs
        ;;
    status)
        check_container_status "nutribox-api"
        check_container_status "nutribox-postgres"
        ;;
    troubleshoot)
        troubleshoot
        ;;
    help)
        show_help
        ;;
    *)
        # Default action if no argument is provided
        if [ -z "$1" ]; then
            check_docker
            start_containers
        else
            error "Unknown command: $1"
            show_help
            exit 1
        fi
        ;;
esac

exit 0