# NutriBox API

NutriBox API is a comprehensive backend solution for a nutrition and meal tracking application. Built with Go and the Fiber framework, it provides robust functionality for user management, meal tracking, recipe management, subscription handling, and more.

## Features

- **User Management**: Registration, authentication, profile management
- **Meal Tracking**: Log and track daily meals with nutritional information
- **Recipe Management**: Store and retrieve recipes
- **Article Management**: Create and manage nutritional articles and categories
- **Health Metrics**: Track user weight, height, and health targets
- **Login Streak**: Track user engagement through login streaks
- **Admin Dashboard**: Manage users, subscriptions, and product tokens

## Tech Stack

- **Language**: Go
- **Framework**: Fiber (Fast HTTP web framework)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT
- **Image Recognition**: LogMeal API integration
- **Documentation**: Swagger

## Prerequisites

- Go 1.16 or higher
- PostgreSQL
- Docker and Docker Compose (optional)

## Installation

### Option 1: Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/nutribox-api.git
cd nutribox-api

# Initialize Go modules
go mod tidy

# Create and configure environment variables
cp .env.example .env
# Edit .env file with your configuration
```

### Option 2: Manual Setup

```bash
# Initialize a new Go project
go mod init <project-name>

# Clone the repository
git clone --depth 1 https://github.com/TheValeHack/nutribox-api.git
cd nutribox-api
rm -rf ./.git

# Install dependencies
go mod tidy

# Set up environment variables
cp .env.example .env
# Edit .env file with your configuration
```

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Server configuration
# Env value : prod || dev
APP_ENV=dev
APP_HOST=0.0.0.0
APP_PORT=3000

# Database configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=nutribox
DB_PORT=5432

# Product token
PRODUCT_TOKEN_EXP_DAYS=30

# LogMeal API
LOG_MEAL_BASE_URL=https://api.logmeal.es/v2
LOG_MEAL_API_KEY=your_logmeal_api_key

# JWT
JWT_SECRET=yoursecretkey
JWT_ACCESS_EXP_MINUTES=30
JWT_REFRESH_EXP_DAYS=30
JWT_RESET_PASSWORD_EXP_MINUTES=10
JWT_VERIFY_EMAIL_EXP_MINUTES=10

# SMTP configuration
SMTP_HOST=email-server
SMTP_PORT=587
SMTP_USERNAME=email-server-username
SMTP_PASSWORD=email-server-password
EMAIL_FROM=support@nutribox.com

# OAuth2 configuration
GOOGLE_CLIENT_ID=yourgoogleclientid.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=yourgoogleclientsecret
REDIRECT_URL=http://localhost:3000/v1/auth/google-callback

## Running the Application

### Local Development

```bash
# Run the application
make start

# Or run with live reload (requires Air to be installed)
air
```

### Using Docker

```bash
# Build and run with Docker Compose
docker-compose up -d

# Or use the provided script
./run-docker.sh
```

## Testing

```bash
# Run all tests
make tests

# Run tests with gotestsum format
make testsum

# Run test for a specific function
make tests-TestUserModel
```

## API Documentation

Swagger documentation is available at `/swagger/index.html` when the application is running.

To generate updated Swagger documentation:

```bash
make swagger
```

## Project Structure

```
src\
 |--config\         # Environment variables and configuration
 |--controller\     # Route controllers (controller layer)
 |--database\       # Database connection & migrations
 |--docs\           # Swagger files
 |--middleware\     # Custom fiber middlewares
 |--model\          # Database models (data layer)
 |--response\       # Response models
 |--router\         # Routes
 |--service\        # Business logic (service layer)
 |--utils\          # Utility classes and functions
 |--validation\     # Request data validation schemas
 |--main.go         # Application entry point
```

### Payment Endpoints

- `POST /v1/subscriptions/purchase/:planID`: Initiate a subscription purchase

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.