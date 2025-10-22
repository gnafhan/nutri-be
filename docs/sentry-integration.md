# Sentry Error Monitoring Integration

This document describes the Sentry error monitoring integration in the NutriBox API.

## Overview

Sentry has been integrated into the NutriBox API to provide comprehensive error monitoring, performance tracking, and debugging capabilities. The integration follows the project's modular architecture patterns.

## Configuration

### Environment Variables

Add the following environment variables to your `.env` file:

```bash
# Sentry Configuration
SENTRY_DSN=https://55de1515c56853c91b984122d47902a5@o4509621033828352.ingest.de.sentry.io/4510219220156496
SENTRY_ENVIRONMENT=development  # or production
SENTRY_DEBUG=true  # Set to false in production
```

### Required Configuration

The `SENTRY_DSN` environment variable is **required** for Sentry to function. If not provided:
- Sentry initialization will be skipped
- A warning will be logged: "SENTRY_DSN not configured, Sentry will not capture errors"
- The application will continue to run without error monitoring

### Default Configuration

If environment variables are not provided, the system will use these defaults:
- **Environment**: `development` (if `APP_ENV=prod`, then `production`)
- **Debug**: `false`

## Architecture

### 1. Configuration Layer (`src/config/config.go`)
- Added Sentry configuration variables
- Integrated with existing Viper configuration system

### 2. Middleware Layer (`src/middleware/sentry.go`)
- `SentryConfig()`: Main Sentry middleware for error capture
- `SentryEnhancement()`: Adds request context and custom tags
- Automatic initialization with proper configuration

### 3. Utility Layer (`src/utils/sentry.go`)
- `SentryService`: Comprehensive utility service for error reporting
- Context-aware error capture with Fiber integration
- Specialized methods for different error types

### 4. Integration Points
- **Main Application** (`src/main.go`): Sentry initialization and middleware setup
- **Router** (`src/router/router.go`): Test endpoints for development

## Features

### Error Monitoring
- Automatic panic recovery and error capture
- Request context preservation
- User context tracking
- Custom tags and metadata

### Performance Monitoring
- Request performance tracking
- Database query monitoring
- External API call tracking

### Security Features
- Sensitive data filtering (Authorization headers, cookies, API keys)
- Before-send hooks for data sanitization
- Environment-based configuration

## Usage Examples

### Basic Error Capture
```go
sentryService := utils.NewSentryService()
err := errors.New("Something went wrong")
sentryService.CaptureErrorWithContext(c, err, nil, nil)
```

### API Error with Context
```go
sentryService := utils.NewSentryService()
sentryService.CaptureAPIError(c, err, "/v1/users", "create_user")
```

### Database Error
```go
sentryService := utils.NewSentryService()
sentryService.CaptureDatabaseError(c, err, "SELECT", "users")
```

### Validation Error
```go
sentryService := utils.NewSentryService()
sentryService.CaptureValidationError(c, err, "email", userEmail)
```

### External API Error
```go
sentryService := utils.NewSentryService()
sentryService.CaptureExternalAPIError(c, err, "LogMeal", "https://api.logmeal.es/v2", 500)
```

### Performance Issue
```go
sentryService := utils.NewSentryService()
sentryService.CapturePerformanceIssue(c, "Slow database query", duration, "user_lookup")
```

### Security Event
```go
sentryService := utils.NewSentryService()
sentryService.CaptureSecurityEvent(c, "Suspicious login attempt", map[string]interface{}{
    "ip_address": c.IP(),
    "user_agent": c.Get("User-Agent"),
})
```

### Adding Breadcrumbs
```go
sentryService := utils.NewSentryService()
sentryService.AddBreadcrumb(c, "User authentication started", "auth", "info", map[string]interface{}{
    "user_id": userID,
})
```

## Testing Endpoints

In development mode, the following test endpoints are available:

- `GET /v1/sentry/test-error` - Tests error reporting
- `GET /v1/sentry/test-message` - Tests message reporting  
- `GET /v1/sentry/test-panic` - Tests panic recovery

## Best Practices

### 1. Error Context
Always provide meaningful context when capturing errors:
```go
tags := map[string]string{
    "operation": "user_creation",
    "endpoint": "/v1/users",
}
extra := map[string]interface{}{
    "user_id": userID,
    "request_data": requestData,
}
sentryService.CaptureErrorWithContext(c, err, tags, extra)
```

### 2. User Context
Set user context for better error tracking:
```go
sentryService.SetUserContext(c, userID, userEmail, username)
```

### 3. Breadcrumbs
Add breadcrumbs for better debugging:
```go
sentryService.AddBreadcrumb(c, "Database query started", "database", "info", map[string]interface{}{
    "query": "SELECT * FROM users WHERE id = ?",
    "params": []interface{}{userID},
})
```

### 4. Performance Monitoring
Monitor slow operations:
```go
start := time.Now()
// ... perform operation ...
duration := time.Since(start)
if duration > 5*time.Second {
    sentryService.CapturePerformanceIssue(c, "Slow operation detected", duration, "user_processing")
}
```

## Security Considerations

1. **Sensitive Data**: The integration automatically filters out sensitive headers and data
2. **Environment Separation**: Different environments are properly tagged
3. **Rate Limiting**: Sentry has built-in rate limiting to prevent spam
4. **Data Retention**: Configure appropriate data retention policies in Sentry

## Monitoring and Alerts

Configure alerts in your Sentry dashboard for:
- Error rate spikes
- Performance degradation
- Security events
- New error types

## Troubleshooting

### Common Issues

1. **Sentry not initializing**: Check DSN configuration and network connectivity
2. **Events not appearing**: Verify environment configuration and debug mode
3. **Performance impact**: Adjust sample rates if needed

### Debug Mode

Enable debug mode to see Sentry initialization and event sending:
```bash
SENTRY_DEBUG=true
```

This will log Sentry operations to help with troubleshooting.

## Dependencies

The integration uses:
- `github.com/getsentry/sentry-go` - Core Sentry SDK
- `github.com/getsentry/sentry-go/fiber` - Fiber integration

These are already included in the project's `go.mod` file.
