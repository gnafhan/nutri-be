# NutriBox API Documentation

## 📚 Documentation Overview

This directory contains comprehensive documentation for the NutriBox API, including the new Freemium Trial System implementation.

## 📁 Documentation Structure

### Core Documentation

- **[`freemium-trial-system.md`](./freemium-trial-system.md)** - Complete guide for frontend developers on the freemium trial system
- **[`frontend-quick-reference.md`](./frontend-quick-reference.md)** - Quick reference guide for frontend developers
- **[`swagger-updates.md`](./swagger-updates.md)** - Swagger/OpenAPI documentation updates

### Additional Resources

- **[`../.cursor/rules/`](../.cursor/rules/)** - Cursor AI rules for development standards
- **[`../Makefile`](../Makefile)** - Build and test commands

## 🚀 Quick Start

### For Frontend Developers

1. **Start Here**: Read [`frontend-quick-reference.md`](./frontend-quick-reference.md) for immediate implementation guidance
2. **Deep Dive**: Review [`freemium-trial-system.md`](./freemium-trial-system.md) for comprehensive understanding
3. **API Reference**: Check Swagger documentation at `/docs` endpoint after running the server

### For Backend Developers

1. **Implementation**: Review the freemium trial system implementation in the codebase
2. **Testing**: Use the test commands in the Makefile (`make test-freemium-no-db`)
3. **Documentation**: Update Swagger docs using [`swagger-updates.md`](./swagger-updates.md)

## 🎯 Freemium Trial System

### What's New

- **2-week free trial** for new users after email verification
- **Full access** to all API features during trial
- **Automatic expiry** after 14 days
- **Seamless upgrade** to paid subscription or product token

### Key Features

- ✅ **Zero friction onboarding** - No payment required for trial
- ✅ **Full feature access** - All endpoints available during trial
- ✅ **Automatic trial creation** - Triggered by email verification
- ✅ **Smart access control** - New middleware handles all access types
- ✅ **Backwards compatible** - Existing users unaffected

## 📋 API Access Levels

| Access Type | Duration | Features | Requirements |
|-------------|----------|----------|--------------|
| **Freemium Trial** | 14 days | Full access | Email verification |
| **Paid Subscription** | Until expiry | Full access | Active subscription |
| **Product Token** | Until expiry | Full access | Valid token |
| **No Access** | - | None | Must upgrade |

## 🔧 Implementation Status

### ✅ Completed

- [x] Freemium trial plan added to database seeder
- [x] Subscription service updated with freemium creation
- [x] Auth service updated to auto-create trial on email verification
- [x] New middleware for freemium access control
- [x] All route files updated to use new middleware
- [x] Comprehensive test suite (unit and integration)
- [x] Test commands added to Makefile
- [x] Frontend documentation created

### 🚀 Ready for Frontend Integration

The backend implementation is **complete and ready** for frontend integration. Frontend developers can now:

1. **Implement registration flow** with trial messaging
2. **Handle email verification** with trial activation
3. **Manage access control** with new error responses
4. **Display trial status** and countdown
5. **Implement upgrade flows** for expired trials

## 🧪 Testing

### Available Test Commands

```bash
# Run freemium tests (no database required)
make test-freemium-no-db

# Run all freemium tests (requires database)
make test-freemium

# Run specific test types
make test-freemium-middleware  # Middleware tests
make test-freemium-service     # Service tests (needs DB)
make test-freemium-integration # Integration tests (needs DB)

# Show all available commands
make test-help
```

### Test Coverage

- ✅ **Middleware Tests** - Access control logic (no database required)
- ✅ **Service Tests** - Subscription creation (requires database)
- ✅ **Integration Tests** - End-to-end flows (requires database)

## 📖 Documentation Usage

### For Frontend Teams

1. **Start with Quick Reference** - Get immediate implementation guidance
2. **Review Full Documentation** - Understand the complete system
3. **Check API Examples** - See real code examples and patterns
4. **Test Integration** - Use the provided test scenarios

### For Backend Teams

1. **Review Implementation** - Check the codebase for implementation details
2. **Run Tests** - Verify functionality with test commands
3. **Update Swagger** - Use the Swagger updates guide
4. **Monitor Deployment** - Track freemium trial metrics

## 🔄 API Changes Summary

### New Behavior

- **Registration**: Users get verification email mentioning 2-week trial
- **Email Verification**: Automatically creates freemium trial with full access
- **Access Control**: New middleware handles freemium, subscription, and product token access
- **Error Responses**: New error codes for trial expiry and access requirements

### Backwards Compatibility

- ✅ **Existing users** - No changes to current access
- ✅ **Existing subscriptions** - Continue to work as before
- ✅ **Existing product tokens** - Continue to work as before
- ✅ **Existing API endpoints** - No breaking changes

## 🚨 Important Notes

### For Frontend Developers

1. **Update Error Handling** - Handle new 403 error responses
2. **Show Trial Status** - Display trial countdown and expiry information
3. **Implement Upgrade Flows** - Guide users to subscription or product token activation
4. **Test Thoroughly** - Verify registration and verification flows

### For Backend Developers

1. **Monitor Metrics** - Track freemium trial creation and conversion rates
2. **Database Seeding** - Ensure Freemium Trial plan is seeded in production
3. **Error Logging** - Monitor 403 responses for access issues
4. **Performance** - Monitor API performance with increased user registration

## 📞 Support

### Documentation Issues

- Check this README for guidance
- Review the specific documentation files
- Contact the backend team for clarification

### Implementation Issues

- Run the test commands to verify functionality
- Check error logs for specific issues
- Review the troubleshooting sections in the documentation

### API Issues

- Verify JWT token validity
- Check subscription status
- Review access control logic

---

**Last Updated**: October 20, 2025  
**Version**: 1.0.0  
**API Version**: v1

For the most up-to-date information, always refer to the latest version of these documentation files.
