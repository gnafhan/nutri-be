# Swagger Documentation Updates for Freemium Trial System

## Overview

This document outlines the necessary updates to the Swagger/OpenAPI documentation to reflect the new freemium trial system implementation.

## Required Updates

### 1. Authentication Endpoints

#### Registration Endpoint (`POST /v1/auth/register`)

**Updated Response Description:**
```yaml
responses:
  201:
    description: "Registration successful. User will receive verification email to start 2-week free trial."
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/RegisterResponse'
        example:
          status: "success"
          message: "Register successfully"
          user:
            id: "123e4567-e89b-12d3-a456-426614174000"
            name: "John Doe"
            email: "john@example.com"
            verified_email: false
          tokens:
            access:
              token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
              expires: "2025-10-20T08:30:00Z"
            refresh:
              token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
              expires: "2025-11-19T08:00:00Z"
```

#### Email Verification Endpoint (`GET /v1/auth/verify-email`)

**Updated Response Description:**
```yaml
responses:
  200:
    description: "Email verified successfully. User now has 2-week free trial with full access to all features."
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/VerifyEmailResponse'
        example:
          status: "success"
          message: "Verify email successfully"
```

### 2. Error Response Schemas

#### New Error Response Types

```yaml
components:
  schemas:
    FreemiumExpiredError:
      type: object
      properties:
        status:
          type: string
          example: "error"
        message:
          type: string
          example: "freemium_expired"
        error:
          type: string
          example: "Your freemium trial has expired. Please activate a product token or purchase a subscription."
        data:
          type: object
          properties:
            upgrade_url:
              type: string
              example: "/v1/subscriptions/plans"

    AccessRequiredError:
      type: object
      properties:
        status:
          type: string
          example: "error"
        message:
          type: string
          example: "access_required"
        error:
          type: string
          example: "Access denied. Please activate a product token or purchase a subscription."
        data:
          type: object
          properties:
            upgrade_url:
              type: string
              example: "/v1/subscriptions/plans"
```

### 3. Protected Endpoints

#### Updated Error Responses for All Protected Endpoints

Add these error responses to all protected endpoints:

```yaml
responses:
  401:
    description: "Unauthorized - Invalid or missing JWT token"
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/UnauthorizedError'
  403:
    description: "Forbidden - No access (freemium expired, no subscription, or no product token)"
    content:
      application/json:
        schema:
          oneOf:
            - $ref: '#/components/schemas/FreemiumExpiredError'
            - $ref: '#/components/schemas/AccessRequiredError'
        examples:
          freemium_expired:
            summary: "Freemium trial expired"
            value:
              status: "error"
              message: "freemium_expired"
              error: "Your freemium trial has expired. Please activate a product token or purchase a subscription."
              data:
                upgrade_url: "/v1/subscriptions/plans"
          access_required:
            summary: "No access"
            value:
              status: "error"
              message: "access_required"
              error: "Access denied. Please activate a product token or purchase a subscription."
              data:
                upgrade_url: "/v1/subscriptions/plans"
```

### 4. User Schema Updates

#### Updated User Schema

```yaml
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          format: uuid
          example: "123e4567-e89b-12d3-a456-426614174000"
        name:
          type: string
          example: "John Doe"
        email:
          type: string
          format: email
          example: "john@example.com"
        verified_email:
          type: boolean
          example: true
          description: "Email verification status. Freemium trial starts after verification."
        role:
          type: string
          example: "user"
        # ... other existing properties
```

### 5. Subscription Schema Updates

#### Updated Subscription Response

```yaml
components:
  schemas:
    UserSubscriptionResponse:
      type: object
      properties:
        id:
          type: string
          format: uuid
        user_id:
          type: string
          format: uuid
        plan:
          $ref: '#/components/schemas/SubscriptionPlanResponse'
        is_active:
          type: boolean
        start_date:
          type: string
          format: date-time
        end_date:
          type: string
          format: date-time
          description: "Trial expires after 14 days from start_date"
        payment_method:
          type: string
          enum: ["freemium_trial", "credit_card", "bank_transfer"]
          example: "freemium_trial"
        payment_status:
          type: string
          enum: ["completed", "pending", "failed"]
          example: "completed"

    SubscriptionPlanResponse:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
          example: "Freemium Trial"
          description: "Plan name. 'Freemium Trial' for 2-week free trial."
        price:
          type: integer
          example: 0
          description: "Price in cents. 0 for freemium trial."
        price_formatted:
          type: string
          example: "Rp 0"
        features:
          type: object
          properties:
            scan_ai:
              type: boolean
              example: true
            scan_calorie:
              type: boolean
              example: true
            chatbot:
              type: boolean
              example: true
            bmi_check:
              type: boolean
              example: true
            weight_tracking:
              type: boolean
              example: true
            health_info:
              type: boolean
              example: true
        validity_days:
          type: integer
          example: 14
          description: "Validity in days. 14 for freemium trial."
        ai_scan_limit:
          type: integer
          example: 10
          description: "AI scan limit per period"
```

### 6. JWT Token Schema Updates

#### Updated Token Claims

```yaml
components:
  schemas:
    JWTClaims:
      type: object
      properties:
        sub:
          type: string
          description: "User ID"
        iat:
          type: integer
          description: "Issued at timestamp"
        exp:
          type: integer
          description: "Expiration timestamp"
        type:
          type: string
          example: "access"
        userData:
          type: object
          properties:
            id:
              type: string
              format: uuid
            name:
              type: string
            email:
              type: string
            role:
              type: string
            verified_email:
              type: boolean
            isProductTokenVerified:
              type: boolean
              description: "Whether user has active product token"
            subscriptionFeatures:
              type: object
              description: "Available features based on subscription or freemium trial"
              properties:
                scan_ai:
                  type: boolean
                scan_calorie:
                  type: boolean
                chatbot:
                  type: boolean
                bmi_check:
                  type: boolean
                weight_tracking:
                  type: boolean
                health_info:
                  type: boolean
```

### 7. Subscription Endpoints

#### Get Subscription Plans

```yaml
/v1/subscriptions/plans:
  get:
    summary: "Get available subscription plans"
    description: "Returns all available subscription plans including Freemium Trial plan"
    responses:
      200:
        description: "List of subscription plans"
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/SuccessWithPaginate[SubscriptionPlanResponse]'
            example:
              status: "success"
              message: "Plans retrieved successfully"
              results:
                - id: "freemium-plan-id"
                  name: "Freemium Trial"
                  price: 0
                  price_formatted: "Rp 0"
                  features:
                    scan_ai: true
                    scan_calorie: true
                    chatbot: true
                    bmi_check: true
                    weight_tracking: true
                    health_info: true
                  validity_days: 14
                  ai_scan_limit: 10
                  is_recommended: false
                  description: "14-day free trial with full features"
                - id: "premium-plan-id"
                  name: "Premium"
                  price: 50000
                  price_formatted: "Rp 50,000"
                  features:
                    scan_ai: true
                    scan_calorie: true
                    chatbot: true
                    bmi_check: true
                    weight_tracking: true
                    health_info: true
                  validity_days: 30
                  ai_scan_limit: 100
                  is_recommended: true
                  description: "Premium subscription with full features"
```

### 8. Access Control Documentation

#### Add Access Control Section

```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: |
        JWT token required for authentication.
        
        Access Levels:
        - **Freemium Trial**: 2-week free trial with full features
        - **Paid Subscription**: Full access until expiry
        - **Product Token**: Full access until expiry
        
        Token includes user subscription information and available features.
```

#### Update Security Requirements

```yaml
security:
  - BearerAuth: []
```

Add this to all protected endpoints with additional description:

```yaml
security:
  - BearerAuth: []
  description: |
    Requires valid JWT token with one of:
    - Active freemium trial (14 days from email verification)
    - Active paid subscription
    - Valid product token
    
    Returns 403 with upgrade information if access is denied.
```

## Implementation Steps

### 1. Update Existing Swagger Files

1. **Update authentication endpoints** with new response descriptions
2. **Add new error response schemas** for freemium and access errors
3. **Update user and subscription schemas** with new fields
4. **Add freemium trial information** to JWT token documentation
5. **Update all protected endpoints** with new error responses

### 2. Regenerate Swagger Documentation

```bash
# Generate updated Swagger documentation
make swagger

# Or manually
cd src && swag init
```

### 3. Verify Documentation

1. **Test authentication flow** in Swagger UI
2. **Verify error responses** match new formats
3. **Check subscription plan responses** include Freemium Trial
4. **Validate JWT token structure** documentation

### 4. Update API Documentation

1. **Add freemium trial section** to main API documentation
2. **Update authentication guide** with new access levels
3. **Add error handling examples** for frontend developers
4. **Include migration guide** for existing integrations

## Example Complete Endpoint Update

```yaml
/v1/meals:
  get:
    summary: "Get user meals"
    description: "Retrieve user's meal history. Requires active subscription, product token, or freemium trial."
    security:
      - BearerAuth: []
    responses:
      200:
        description: "Meals retrieved successfully"
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/SuccessWithPaginate[Meal]'
      401:
        description: "Unauthorized - Invalid or missing JWT token"
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UnauthorizedError'
      403:
        description: "Forbidden - No access or freemium trial expired"
        content:
          application/json:
            schema:
              oneOf:
                - $ref: '#/components/schemas/FreemiumExpiredError'
                - $ref: '#/components/schemas/AccessRequiredError'
            examples:
              freemium_expired:
                summary: "Freemium trial expired"
                value:
                  status: "error"
                  message: "freemium_expired"
                  error: "Your freemium trial has expired. Please activate a product token or purchase a subscription."
                  data:
                    upgrade_url: "/v1/subscriptions/plans"
              access_required:
                summary: "No access"
                value:
                  status: "error"
                  message: "access_required"
                  error: "Access denied. Please activate a product token or purchase a subscription."
                  data:
                    upgrade_url: "/v1/subscriptions/plans"
```

## Testing Swagger Updates

### 1. Test Authentication Flow

1. **Register new user** via Swagger UI
2. **Verify email** (simulate with test token)
3. **Login and check JWT** contains subscription features
4. **Access protected endpoint** with freemium access

### 2. Test Error Responses

1. **Test 401 responses** with invalid/missing JWT
2. **Test 403 responses** with expired freemium
3. **Test 403 responses** with no access
4. **Verify error formats** match documentation

### 3. Test Subscription Endpoints

1. **Get subscription plans** and verify Freemium Trial is listed
2. **Check plan features** are correctly documented
3. **Verify pricing information** is accurate

---

**Note**: After implementing these Swagger updates, regenerate the documentation using `make swagger` and verify all changes are reflected in the Swagger UI.
