# Security and Authentication

## Overview

Golf Gamez is designed as an anonymous, shareable application where users can create games without registration. Security focuses on preventing abuse while maintaining ease of use for weekend golfers.

## Authentication Model

### Anonymous Access
- **No user registration required**
- **No authentication for basic game creation**
- Games are identified by shareable tokens
- Anyone with a game token can view and participate

### Token-Based Access Control

#### Game Access Tokens
```typescript
interface GameTokens {
  share_token: string;      // Full access - can modify game and scores
  spectator_token: string;  // Read-only access - can view but not modify
}
```

**Share Token (`share_token`)**
- Provides full read/write access to the game
- Used by players to join and record scores
- Should be shared only with active participants
- Format: `gt_abc123def456ghi789` (30 characters)

**Spectator Token (`spectator_token`)**
- Provides read-only access to view live scores
- Safe to share publicly or with non-participants
- Cannot modify game state or scores
- Format: `st_abc123def456ghi789` (30 characters)

### Token Validation

#### Share Token Access
```http
GET /api/games/{shareToken}
Authorization: Bearer gt_abc123def456ghi789

POST /api/games/{shareToken}/players
Authorization: Bearer gt_abc123def456ghi789
```

#### Spectator Token Access
```http
GET /api/spectate/{spectatorToken}
GET /api/spectate/{spectatorToken}/leaderboard
```

## Security Measures

### Rate Limiting

#### Per-IP Limits
- **Game Creation**: 10 games per hour per IP
- **Score Updates**: 100 requests per minute per IP
- **API Calls**: 1000 requests per hour per IP

#### Per-Game Limits
- **Players**: Maximum 4 players per game
- **Score Updates**: 1 score update per hole per player per minute
- **Game Duration**: Auto-abandon games older than 24 hours without activity

### Input Validation

#### String Sanitization
- Player names: HTML encoded, max 100 characters
- No script tags or executable content allowed
- Unicode normalization applied

#### Numeric Validation
- Scores: 1-20 strokes per hole (reasonable golf limits)
- Putts: 0-10 putts per hole (0 for hole-in-one, max 10 reasonable)
- Handicaps: 0-54 (official USGA range)
- Holes: 1-18 only

#### Business Logic Validation
- Cannot record scores for future holes
- Cannot modify scores after game completion
- Cannot add players after game starts
- Score consistency checks (putts ≤ strokes)

### Data Protection

#### Token Security
- Tokens generated using cryptographically secure random number generation
- Tokens are single-use for creation, persistent for access
- No token recycling - each game gets unique tokens

#### Data Retention
- Games auto-deleted after 30 days of inactivity
- No personal information stored beyond game context
- Anonymous usage - no PII collection

#### Database Security
- Prepared statements prevent SQL injection
- Input sanitization at API layer
- Foreign key constraints maintain data integrity

### Cross-Origin Resource Sharing (CORS)

```http
Access-Control-Allow-Origin: https://golfgamez.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Max-Age: 86400
```

**Allowed Origins:**
- Production: `https://golfgamez.com`
- Development: `http://localhost:3000`, `http://localhost:8080`
- Mobile: Custom scheme `golfgamez://`

### Content Security Policy (CSP)

```http
Content-Security-Policy:
  default-src 'self';
  script-src 'self' 'unsafe-inline';
  style-src 'self' 'unsafe-inline';
  img-src 'self' data:;
  connect-src 'self' wss:;
  font-src 'self';
  object-src 'none';
  base-uri 'self';
  form-action 'self';
```

### API Abuse Prevention

#### Duplicate Game Detection
- Prevent creation of games with identical player names within 1 hour
- IP-based cooldown for rapid game creation

#### Score Manipulation Prevention
- Timestamps on all score entries
- Score history tracking (updates vs creates)
- Reasonable score limits based on golf norms

#### Bot Protection
- Simple CAPTCHA for rapid game creation
- Behavioral analysis for non-human patterns
- Temporary IP blocking for abuse

### WebSocket Security

#### Connection Authentication
```javascript
// WebSocket connection with token
const ws = new WebSocket('wss://api.golfgamez.com/ws/games/gt_abc123def456');

// Token sent in first message
ws.send(JSON.stringify({
  type: 'auth',
  token: 'gt_abc123def456ghi789'
}));
```

#### Message Validation
- All incoming messages validated against schema
- Rate limiting applied to WebSocket messages
- Automatic disconnection for invalid tokens

### Error Handling Security

#### Information Disclosure Prevention
- Generic error messages for failed authentication
- No internal system information in error responses
- Consistent response times for valid/invalid tokens

#### Error Response Examples

**Valid Token, Game Not Found:**
```json
{
  "error": {
    "code": "game_not_found",
    "message": "Game not found"
  }
}
```

**Invalid Token:**
```json
{
  "error": {
    "code": "invalid_token",
    "message": "Invalid or expired token"
  }
}
```

### Monitoring and Logging

#### Security Event Logging
- Failed token validation attempts
- Rapid request patterns
- Unusual score entries (e.g., impossible scores)
- Game creation spikes from single IPs

#### Metrics Tracking
- Request latency by endpoint
- Token validation success rates
- Game completion rates
- Error rate monitoring

### Privacy Considerations

#### Data Minimization
- Only collect necessary game data
- No tracking cookies or analytics
- No email or phone number collection

#### Anonymous Usage
- No user accounts or profiles
- Games identified only by tokens
- No cross-game player identification

#### Data Sharing
- No third-party data sharing
- No advertising or tracking pixels
- Self-contained application

## Threat Model

### Identified Threats

#### Low Risk
- **Score Cheating**: Players can manipulate their own scores
  - *Mitigation*: Social enforcement, reasonable limits
- **Game Viewing**: Spectators can view any game with token
  - *Mitigation*: Separate spectator tokens, expected behavior

#### Medium Risk
- **Token Sharing**: Share tokens could be leaked publicly
  - *Mitigation*: Game auto-deletion, limited game duration
- **API Abuse**: Automated game creation or score spam
  - *Mitigation*: Rate limiting, behavior detection

#### High Risk
- **Data Injection**: SQL injection or XSS attacks
  - *Mitigation*: Prepared statements, input sanitization, CSP
- **Service Disruption**: DDoS or resource exhaustion
  - *Mitigation*: Rate limiting, load balancing, monitoring

### Security Testing

#### Automated Testing
- Input validation testing with malformed data
- Rate limiting verification
- Token security testing

#### Manual Testing
- Penetration testing for common vulnerabilities
- Social engineering resistance
- Token lifecycle testing

## Compliance

### No PII Collection
- Application does not collect personally identifiable information
- Player names are not linked to real identities
- No GDPR or CCPA compliance requirements

### Industry Standards
- Follows OWASP guidelines for web application security
- Uses TLS 1.3 for all communications
- Implements security headers per Mozilla guidelines

## Future Security Enhancements

### Optional User Accounts
- Future feature for persistent game history
- OAuth integration for social login
- Enhanced security for registered users

### Advanced Analytics
- Machine learning for bot detection
- Behavioral analysis for score validation
- Predictive abuse prevention

### Enhanced Monitoring
- Real-time threat detection
- Automated incident response
- Security dashboard for operators