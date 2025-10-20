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
const hasAccess = user.subscriptionFeatures?.scan_ai === true;
// Freemium users have full features for 14 days
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

### Product Token
- ✅ **Full access** to all features
- ✅ **All API endpoints** available
- ⏰ **Until expiry** date
- 🎫 **Token-based** access

### No Access
- ❌ **403 Forbidden** responses
- 📝 **Upgrade prompts** required
- 💳 **Must subscribe** or activate token

## 🚨 Error Codes

| Code | Message | Description | Action |
|------|---------|-------------|--------|
| `401` | `Please authenticate` | Missing/invalid JWT | Redirect to login |
| `403` | `freemium_expired` | Trial ended | Show upgrade modal |
| `403` | `access_required` | No subscription/token | Show upgrade modal |

## 💡 Frontend Implementation Examples

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
