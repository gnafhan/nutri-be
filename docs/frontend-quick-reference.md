# Frontend Developer Quick Reference - Freemium Trial System

## 🚀 Quick Start

### New User Registration Flow
```javascript
// 1. Register user
const response = await api.post('/v1/auth/register', userData);
// → User gets verification email

// 2. User verifies email
// → Freemium trial automatically created (14 days)

// 3. User can immediately access all features
```

### Access Control Check
```javascript
// Check subscription type
const isFreemium = user.subscriptionType === 'freemium';
const isActive = user.subscriptionStatus === 'active';

// Check specific features
const hasScanAccess = user.subscriptionFeatures?.scan_ai === true;
const hasCalorieAccess = user.subscriptionFeatures?.scan_calorie === true;

// Freemium users have full features for 14 days
if (isFreemium && isActive) {
  // User has full access to all features
}
```

### Error Handling
```javascript
// Handle freemium expiry
if (error.response?.data?.message === 'freemium_expired') {
  showUpgradeModal();
}
```

## 📋 API Endpoints

| Endpoint | Method | Description | Access Required |
|----------|--------|-------------|-----------------|
| `/v1/auth/register` | POST | Register new user | None |
| `/v1/auth/verify-email` | GET | Verify email + create freemium trial | None |
| `/v1/auth/login` | POST | Login user | None |
| `/v1/subscriptions/plans` | GET | Get subscription plans | None |
| `/v1/auth/verify-product-token` | POST | Activate product token | JWT |
| `/v1/meals` | GET | Get user meals | JWT + Access |
| `/v1/meals/scan` | POST | AI meal scanning | JWT + Access |
| `/v1/product-token/verify` | POST | Verify product token (coupon) | JWT |

## 🔐 Access Levels

### Freemium Trial (14 days)
- ✅ **Full access** to all features
- ✅ **All API endpoints** available
- ⏰ **14 days** from email verification
- 🔄 **Auto-expires** after trial period

### Paid Subscription
- ✅ **Full access** to all features
- ✅ **All API endpoints** available
- ⏰ **Until expiry** date
- 💳 **Paid** subscription

### Product Token (Coupon)
- ✅ **Full access** to all features
- ✅ **All API endpoints** available
- ⏰ **Until expiry** date
- 🎫 **Token-based** access
- 🔄 **Upgrades freemium** to paid subscription

### No Access
- ❌ **403 Forbidden** responses
- 📝 **Upgrade prompts** required
- 💳 **Must subscribe** or activate token

## 🔑 JWT Payload Structure

The JWT token now includes subscription information directly in the payload:

```javascript
{
  "sub": "user-uuid",
  "iat": 1234567890,
  "exp": 1234567890,
  "type": "access",
  "userData": {
    "id": "user-uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "verified_email": true,
    "subscriptionType": "freemium",        // "freemium", "subscription", or "none"
    "subscriptionStatus": "active",        // "active" or "inactive"
    "subscriptionFeatures": {
      "scan_ai": true,
      "scan_calorie": true,
      "chatbot": true,
      "bmi_check": true,
      "weight_tracking": true,
      "health_info": true
    },
    "subscriptionStartDate": "2025-10-20T22:12:40+07:00",  // ISO 8601 format or null
    "subscriptionEndDate": "2025-11-03T22:12:40+07:00",    // ISO 8601 format or null
    "isProductTokenVerified": false
  }
}
```

### Subscription Types
- `"freemium"` - 14-day free trial with full features
- `"subscription"` - Any paid plan (Hemat, Early Bird, Sehat, Sultan)
- `"none"` - No active subscription

### Subscription Date Fields
- `subscriptionStartDate` - When the subscription started (ISO 8601 format or null)
- `subscriptionEndDate` - When the subscription expires (ISO 8601 format or null)

## 🎫 Product Token (Coupon) System

Product tokens allow users to claim subscription plans using coupon codes. This is useful for:
- **Promotional campaigns**
- **Gift subscriptions**
- **Special offers**
- **Upgrading freemium users**

### How Product Tokens Work

1. **Admin creates product tokens** with associated subscription plans
2. **Users verify tokens** via `POST /v1/product-token/verify?token=COUPON_CODE`
3. **System automatically**:
   - Deactivates any active freemium subscription
   - Creates new subscription based on token's plan
   - Sets payment status to "success"
   - Links subscription to the product token

### Frontend Implementation

```javascript
// Verify product token
const verifyProductToken = async (token) => {
  try {
    const response = await fetch('/v1/product-token/verify?token=' + token, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.ok) {
      // Token verified successfully
      // User's subscription will be updated
      // Refresh user data to get new subscription info
      await refreshUserData();
      showSuccessMessage('Coupon applied successfully!');
    } else {
      const error = await response.json();
      showErrorMessage(error.message);
    }
  } catch (error) {
    showErrorMessage('Failed to verify coupon');
  }
};

// Usage example
const handleCouponSubmit = (couponCode) => {
  verifyProductToken(couponCode);
};
```

### Product Token Response

```javascript
// Success Response
{
  "status": "success",
  "message": "Verify product token successfully"
}

// Error Responses
{
  "status": "error",
  "message": "Invalid or already used product token"  // 404
}

{
  "status": "error", 
  "message": "Can only be connected with 1 product token."  // 403
}
```

## 🚨 Error Codes

| Code | Message | Description | Action |
|------|---------|-------------|--------|
| `401` | `Please authenticate` | Missing/invalid JWT | Redirect to login |
| `403` | `freemium_expired` | Trial ended | Show upgrade modal |
| `403` | `access_required` | No subscription/token | Show upgrade modal |

## 💡 Frontend Implementation Examples

### React/JavaScript Example
```javascript
// Decode JWT token (you'll need a JWT library like jwt-decode)
import jwtDecode from 'jwt-decode';

const token = localStorage.getItem('access_token');
const decoded = jwtDecode(token);
const user = decoded.userData;

// Check subscription status
const isFreemium = user.subscriptionType === 'freemium';
const isPaidSubscription = user.subscriptionType === 'subscription';
const isActive = user.subscriptionStatus === 'active';
const hasFullAccess = (isFreemium || isPaidSubscription) && isActive;

// Check specific features
const canScanAI = user.subscriptionFeatures?.scan_ai === true;
const canUseChatbot = user.subscriptionFeatures?.chatbot === true;

// Check subscription dates
const startDate = user.subscriptionStartDate ? new Date(user.subscriptionStartDate) : null;
const endDate = user.subscriptionEndDate ? new Date(user.subscriptionEndDate) : null;
const daysRemaining = endDate ? Math.ceil((endDate - new Date()) / (1000 * 60 * 60 * 24)) : 0;

// Show/hide features based on subscription
if (canScanAI) {
  // Show AI scanning feature
  showAIScanButton();
} else {
  // Show upgrade prompt
  showUpgradeModal();
}

// Display subscription info
if (isFreemium) {
  showTrialBanner(`You're on a 14-day free trial (${daysRemaining} days remaining)`);
} else if (isPaidSubscription) {
  showSubscriptionInfo(`Active subscription (${daysRemaining} days remaining)`);
} else {
  showUpgradePrompt();
}
```

### Vue.js Example
```javascript
// In your Vue component
computed: {
  userSubscription() {
    const token = this.$store.getters['auth/accessToken'];
    if (!token) return null;
    
    const decoded = jwtDecode(token);
    return decoded.userData;
  },
  
  isFreemium() {
    return this.userSubscription?.subscriptionType === 'freemium';
  },
  
  hasFeatureAccess() {
    return (feature) => {
      return this.userSubscription?.subscriptionFeatures?.[feature] === true;
    };
  }
}
```

### User Status Component
```javascript
const UserStatus = ({ user }) => {
  const features = user.subscriptionFeatures || {};
  
  if (features.scan_ai) {
    // Has access (freemium, subscription, or product token)
    const daysLeft = calculateDaysLeft(user.trialEndDate);
    return (
      <div>
        <h3>✅ Full Access Active</h3>
        {daysLeft > 0 && <p>⏰ {daysLeft} days left in trial</p>}
      </div>
    );
  }
  
  // No access - show upgrade prompt
  return (
    <div>
      <h3>❌ Upgrade Required</h3>
      <button onClick={showUpgradeModal}>Upgrade Now</button>
    </div>
  );
};
```

### API Error Handler
```javascript
const handleApiError = (error) => {
  const { status, data } = error.response;
  
  switch (data.message) {
    case 'freemium_expired':
      showUpgradeModal({
        title: 'Trial Expired',
        message: 'Your 2-week trial has ended. Upgrade to continue.',
        type: 'trial_expired'
      });
      break;
      
    case 'access_required':
      showUpgradeModal({
        title: 'Upgrade Required',
        message: 'Activate a product token or subscribe to access NutriBox.',
        type: 'no_access'
      });
      break;
      
    case 'Please authenticate':
      redirectToLogin();
      break;
  }
};
```

### Trial Countdown Component
```javascript
const TrialCountdown = ({ user }) => {
  const [timeLeft, setTimeLeft] = useState('');
  
  useEffect(() => {
    const updateCountdown = () => {
      const now = new Date();
      const endDate = new Date(user.trialEndDate);
      const diff = endDate - now;
      
      if (diff <= 0) {
        setTimeLeft('Trial Expired');
        return;
      }
      
      const days = Math.floor(diff / (1000 * 60 * 60 * 24));
      const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
      
      setTimeLeft(`${days}d ${hours}h left`);
    };
    
    updateCountdown();
    const interval = setInterval(updateCountdown, 1000 * 60 * 60); // Update hourly
    
    return () => clearInterval(interval);
  }, [user.trialEndDate]);
  
  return (
    <div className="trial-countdown">
      <span className="trial-badge">FREE TRIAL</span>
      <span className="time-left">{timeLeft}</span>
    </div>
  );
};
```

## 🧪 Testing Checklist

### Registration Flow
- [ ] User can register successfully
- [ ] Verification email is mentioned
- [ ] Success message mentions 2-week trial

### Email Verification
- [ ] Email verification works
- [ ] Success message mentions trial activation
- [ ] User gets immediate access after verification

### Access Control
- [ ] Freemium users can access all endpoints
- [ ] Trial countdown displays correctly
- [ ] Expired trial shows upgrade prompt

### Error Handling
- [ ] 401 errors redirect to login
- [ ] 403 freemium_expired shows upgrade modal
- [ ] 403 access_required shows upgrade modal

### Upgrade Flow
- [ ] Subscription plans are displayed
- [ ] Product token activation works
- [ ] Payment flow completes successfully

## 🔧 Configuration

### Environment Variables
```javascript
// API base URL
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Feature flags
const FEATURES = {
  FREEMIUM_TRIAL: true,
  SUBSCRIPTION_PLANS: true,
  PRODUCT_TOKENS: true
};
```

### API Configuration
```javascript
// Axios configuration
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Add JWT token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('jwt_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

## 📱 Mobile Considerations

### Deep Linking
```javascript
// Handle email verification deep links
const handleDeepLink = (url) => {
  if (url.includes('verify-email')) {
    const token = extractTokenFromUrl(url);
    verifyEmail(token);
  }
};
```

### Offline Handling
```javascript
// Cache trial status for offline use
const cacheTrialStatus = (user) => {
  localStorage.setItem('trial_status', JSON.stringify({
    hasAccess: user.subscriptionFeatures?.scan_ai === true,
    trialEndDate: user.trialEndDate,
    cachedAt: new Date().toISOString()
  }));
};
```

## 🎨 UI/UX Guidelines

### Trial Status Indicators
- 🟢 **Green badge**: Active trial/subscription
- 🟡 **Yellow badge**: Trial expiring soon (< 3 days)
- 🔴 **Red badge**: Trial expired
- ⚫ **Gray badge**: No access

### Upgrade Prompts
- **Prominent placement**: Top of screen or modal
- **Clear messaging**: Explain benefits of upgrade
- **Easy access**: Direct links to subscription plans
- **Urgency indicators**: For expiring trials

### Success Messages
- **Registration**: "Check your email to verify and start your 2-week free trial"
- **Verification**: "Email verified! Your 2-week free trial is now active"
- **Upgrade**: "Welcome to Premium! Enjoy unlimited access"

## 🚀 Deployment Checklist

### Pre-deployment
- [ ] Test registration flow end-to-end
- [ ] Verify error handling works correctly
- [ ] Test upgrade flows complete successfully
- [ ] Check trial countdown accuracy

### Post-deployment
- [ ] Monitor error rates for 403 responses
- [ ] Track freemium trial conversion rates
- [ ] Monitor API performance with new users
- [ ] Verify email delivery for verification

## 📞 Support

### Common Issues
1. **"Please authenticate" errors**: Check JWT token storage
2. **Trial not starting**: Verify email verification completed
3. **Access denied**: Check subscription status and features
4. **Upgrade not working**: Verify payment flow completion

### Debug Information
```javascript
// Log user access level for debugging
console.log('User Access Level:', {
  hasAccess: user.subscriptionFeatures?.scan_ai === true,
  accessType: getUserAccessType(user),
  trialEndDate: user.trialEndDate,
  features: user.subscriptionFeatures
});
```

---

**Need Help?** Check the full documentation in `docs/freemium-trial-system.md` or contact the backend team.
