## Quick Start

To create a project, simply run:

```bash
go mod init <project-name>
```

## Manual Installation

If you would still prefer to do the installation manually, follow these steps:

Clone the repo:

```bash
git clone --depth 1 https://github.com/TheValeHack/nutribox-api.git
cd nutribox-api
rm -rf ./.git
```

Install the dependencies:

```bash
go mod tidy
```

Set the environment variables:

```bash
cp .env.example .env

# open .env and modify the environment variables (if needed)
```

## Commands

Running locally:

```bash
make start
```

Or running with live reload:

```bash
air
```

Testing:

```bash
# run all tests
make tests

# run all tests with gotestsum format
make testsum

# run test for the selected function name
make tests-TestUserModel
```

Linting:

```bash
# run lint
make lint
```

Swagger:

```bash
# generate the swagger documentation
make swagger
```

## Environment Variables

The environment variables can be found and modified in the `.env` file. They come with these default values:

```bash
# server configuration
# Env value : prod || dev
APP_ENV=dev
APP_HOST=0.0.0.0
APP_PORT=3000

# database configuration
DB_HOST=postgresdb
DB_USER=postgres
DB_PASSWORD=thisisasamplepassword
DB_NAME=fiberdb
DB_PORT=5432

# JWT
# JWT secret key
JWT_SECRET=thisisasamplesecret
# Number of minutes after which an access token expires
JWT_ACCESS_EXP_MINUTES=30
# Number of days after which a refresh token expires
JWT_REFRESH_EXP_DAYS=30
# Number of minutes after which a reset password token expires
JWT_RESET_PASSWORD_EXP_MINUTES=10
# Number of minutes after which a verify email token expires
JWT_VERIFY_EMAIL_EXP_MINUTES=10

# SMTP configuration options for the email service
SMTP_HOST=email-server
SMTP_PORT=587
SMTP_USERNAME=email-server-username
SMTP_PASSWORD=email-server-password
EMAIL_FROM=support@yourapp.com

# OAuth2 configuration
GOOGLE_CLIENT_ID=yourapps.googleusercontent.com
GOOGLE_CLIENT_SECRET=thisisasamplesecret
REDIRECT_URL=http://localhost:3000/v1/auth/google-callback
```

## Project Structure

```
src\
 |--config\         # Environment variables and configuration related things
 |--controller\     # Route controllers (controller layer)
 |--database\       # Database connection & migrations
 |--docs\           # Swagger files
 |--middleware\     # Custom fiber middlewares
 |--model\          # Postgres models (data layer)
 |--response\       # Response models
 |--router\         # Routes
 |--service\        # Business logic (service layer)
 |--utils\          # Utility classes and functions
 |--validation\     # Request data validation schemas
 |--main.go         # Fiber app
```

# Subscription with Midtrans Integration

This API supports payment processing with Midtrans for subscription plans. 

## Setup

1. Create a Midtrans account at https://midtrans.com
2. Get your Server Key and Client Key from the Midtrans Dashboard
3. Add the following environment variables to your `.env` file:

```
MIDTRANS_MERCHANT_ID=your_midtrans_merchant_id
NEXT_PUBLIC_MIDTRANS_CLIENT_KEY=your_midtrans_client_key
MIDTRANS_SERVER_KEY=your_midtrans_server_key
MIDTRANS_STATUS=SANDBOX or PRODUCTION
```

## API Endpoints

### Purchase Subscription

`POST /subscriptions/purchase/:planID`

Request body (optional):
```json
{
  "payment_method": "credit_card | gopay | shopeepay | bank_transfer"
}
```

If no payment method is specified, Midtrans will display all available payment options on the payment page.

Response:
```json
{
  "status": "success",
  "message": "Payment initiated successfully",
  "data": {
    "transaction_token": "snap-token-from-midtrans",
    "redirect_url": "https://app.midtrans.com/snap/v2/vtweb/...",
    "order_id": "SUB-12345-1614849849"
  }
}
```

### Payment Notification Webhook

Set up a webhook URL in the Midtrans Dashboard to receive payment notifications:

1. Login to your Midtrans Dashboard
2. Go to Settings > Configuration
3. Set the Payment Notification URL to `https://your-api-domain.com/v1/subscriptions/notification`

## Frontend Integration

To display the Snap payment page:

1. Include the Snap.js library in your HTML:
```html
<script src="https://app.sandbox.midtrans.com/snap/snap.js" data-client-key="YOUR-CLIENT-KEY"></script>
```

2. Use the transaction token to display the payment page:
```javascript
snap.pay('TRANSACTION_TOKEN', {
  onSuccess: function(result){
    // Handle success
  },
  onPending: function(result){
    // Handle pending
  },
  onError: function(result){
    // Handle error
  }
});
```