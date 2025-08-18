## Gateway Service

Go-based HTTP API Gateway that fronts the microservices via gRPC and exposes REST endpoints, Swagger docs, metrics, and public web pages (payment).

### Features

- **HTTP Gateway**: REST API over gRPC microservices (Auth, User, Chat, Payment)
- **Auth & RBAC**: JWT-based authentication and role checks (User, PremiumUser, Admin)
- **Swagger**: Interactive API docs at `/swagger`
- **Metrics**: Prometheus metrics at `/metrics`
- **WebSocket**: Chat WebSocket endpoint
- **Payment Pages**: Server-rendered payment pages (checkout, success, failed)
- **Config-Driven**: YAML/env config using Viper
- **Structured Logging**: Zap-based logging

## Quick Start

### Prerequisites
- Go 1.23+
- Dependent microservices running (or configure their ports):
  - Auth, User, Chat, Payment gRPC services

### Run Service
```bash
cd services/gateway
go mod download
go run cmd/main.go
```

Service runs on port `8080` by default.

## Configuration

Create `services/gateway/internal/config/config.yaml` or set environment variables.

```yaml
environment: development

http:
  port: "8080"
  read_timeout: 10
  write_timeout: 10
  idle_timeout: 60

services:
  auth_port: "50051"
  admin_port: "50052"
  user_port: "50053"
  chat_port: "50054"
  payment_port: "50055"

auth:
  jwt:
    secret_key: "your-secret-key"
    access_token_minutes: 15
    refresh_token_days: 7
    issuer: "qubool-kallyanam"

otp:
  length: 6

base_url: "http://localhost:8080"
default_timezone: "Asia/Kolkata"
```

Config file path can be overridden via `CONFIG_PATH` env var.

## API Endpoints

Base path: `/api/v1`

### Auth
- User
  - `POST /auth/user/register` – Register
  - `POST /auth/user/verify` – Verify OTP
  - `POST /auth/user/login` – Login
  - `POST /auth/user/logout` – Logout (JWT + User role)
  - `POST /auth/user/delete` – Delete account (JWT + User role)
  - `POST /auth/user/refresh` – Refresh access token
- Admin
  - `POST /auth/admin/login` – Login
  - `POST /auth/admin/logout` – Logout (JWT + Admin role)
  - `POST /auth/admin/block-user` – Block user (JWT + Admin role)
  - `POST /auth/admin/unblock-user` – Unblock user (JWT + Admin role)
  - `GET /auth/admin/users` – List users (JWT + Admin role)
  - `GET /auth/admin/user` – Get user by field (JWT + Admin role)

### User
- Profile
  - `PATCH /user/profile` – Partial update (JWT + User)
  - `PUT /user/profile` – Replace (JWT + User)
  - `GET /user/profile` – Get (JWT + User)
  - `GET /user/profiles/:profile_id` – Get full details for user (JWT + User)
  - `GET /user/profile-details/:profile_id` – Get full details for admin (JWT + Admin)
- Photos
  - `POST /user/profile/profile-photo` – Get upload URL (JWT + User)
  - `POST /user/profile/profile-photo/confirm` – Confirm upload (JWT + User)
  - `DELETE /user/profile/profile-photo` – Delete (JWT + User)
  - `POST /user/profile/additional-photo` – Get upload URL (JWT + User)
  - `POST /user/profile/additional-photo/confirm` – Confirm upload (JWT + User)
  - `DELETE /user/profile/additional-photo/:display_order` – Delete (JWT + User)
  - `GET /user/profile/additional-photos` – List (JWT + User)
- Partner Preference
  - `POST /user/preference` – Create (JWT + User)
  - `PATCH /user/preference` – Update (JWT + User)
  - `GET /user/preference` – Get (JWT + User)
- Matching
  - `GET /user/recommendations` – Recommendations (JWT + User)
  - `POST /user/match-action` – Record action (JWT + User)
  - `PUT /user/match-action` – Update action (JWT + User)
  - `GET /user/matches/liked` – Liked profiles (JWT + User)
  - `GET /user/matches/passed` – Passed profiles (JWT + User)
  - `GET /user/matches/mutual` – Mutual matches (JWT + User)

### Chat
- `POST /chat/conversation` – Create conversation (JWT + PremiumUser)
- `GET /chat/conversation/:conversation_id/messages` – Get messages (JWT + PremiumUser)
- `GET /chat/ws` – WebSocket endpoint

### Payment (API)
- `GET /payment/subscription-plans` – List active plans
- `GET /payment/subscription-plan` – Get plan by ID
- `POST /payment/subscription-plan` – Create (JWT + Admin)
- `PATCH /payment/subscription-plan` – Update (JWT + Admin)
- `POST /payment/order` – Create order (JWT + User)
- `GET /payment/subscriptions` – Active subscription (JWT + User)
- `GET /payment/payments-history` – Payment history (JWT + User)
- `GET /payment/admin/completed-payments` – Completed payments (JWT + Admin)

### Payment (Web Pages)
- `GET /payment/checkout`
- `GET /payment/verify`
- `GET /payment/success`
- `GET /payment/failed`

### Docs & Metrics
- `GET /swagger/*any` – Swagger UI (e.g., `/swagger/index.html`)
- `GET /metrics` – Prometheus metrics

## Environment Variables

```bash
# config
export CONFIG_PATH=./internal/config/config.yaml
export ENVIRONMENT=development

# http
export HTTP_PORT=8080
export HTTP_READ_TIMEOUT=10
export HTTP_WRITE_TIMEOUT=10
export HTTP_IDLE_TIMEOUT=60

# services (gRPC targets)
export SERVICES_AUTH_PORT=50051
export SERVICES_ADMIN_PORT=50052
export SERVICES_USER_PORT=50053
export SERVICES_CHAT_PORT=50054
export SERVICES_PAYMENT_PORT=50055

# auth/jwt
export AUTH_JWT_SECRET_KEY=your-secret-key
export AUTH_JWT_ACCESS_TOKEN_MINUTES=15
export AUTH_JWT_REFRESH_TOKEN_DAYS=7
export AUTH_JWT_ISSUER=qubool-kallyanam

# misc
export OTP_LENGTH=6
export BASE_URL=http://localhost:8080
export DEFAULT_TIMEZONE=Asia/Kolkata
```

Note: keys map to Viper with dot-to-underscore conversion (e.g., `auth.jwt.secret_key` -> `AUTH_JWT_SECRET_KEY`).

## Dependencies

- `github.com/gin-gonic/gin` – HTTP framework
- `go.uber.org/zap` – Logging
- `github.com/spf13/viper` – Configuration
- `github.com/swaggo/gin-swagger` / `github.com/swaggo/files` – Swagger UI
- `github.com/prometheus/client_golang` – Metrics exporter
- `google.golang.org/grpc` – gRPC clients
- `github.com/golang-jwt/jwt/v5` – JWT manager (via shared `pkg/security/jwt`)
- Internal clients/usecases/handlers for Auth, User, Chat, Payment

## Notes

- Templates for payment pages: `internal/web/templates`
- Swagger docs base path: `/api/v1`
- Production: enable TLS for gRPC clients and set `environment=production` to use Gin release mode.