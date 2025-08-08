# Auth Service

Go-based authentication microservice with gRPC API for user registration, login, and admin operations.

## Features

- **User Auth**: Registration with OTP, login/logout, account deletion
- **Admin Auth**: Admin login, user management, user blocking
- **JWT Tokens**: Access & refresh token management
- **Security**: Bcrypt passwords, role-based access control

## Quick Start

### Prerequisites
- Go 1.23+
- PostgreSQL
- Redis

### Run Service
```bash
cd services/auth
go mod download
go run cmd/main.go
```

Service runs on port `50051` by default.

### Database Setup
```bash
# Run migrations in order:
psql -d your_db -f migrations/postgres/20250504120153_create_pending_registrations.up.sql
psql -d your_db -f migrations/postgres/20250504120523_create_users.up.sql
psql -d your_db -f migrations/postgres/20250506082116_create_admins.up.sql
```

## Configuration

Set environment variables or create `config/config.yaml`:

```yaml
grpc:
  port: 50051

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: qubool_kallyanam_auth

redis:
  host: localhost
  port: 6379

auth:
  jwt:
    secret_key: "your-secret-key"
    access_token_minutes: 15
    refresh_token_days: 7
```

## API Endpoints

### User Operations
- `UserRegister` - Register new user
- `UserVerification` - Verify registration with OTP
- `UserLogin` - User login
- `UserLogout` - User logout
- `UserDelete` - Delete user account
- `RefreshToken` - Refresh access token

### Admin Operations  
- `AdminLogin` - Admin login
- `AdminLogout` - Admin logout
- `BlockUser` - Block user by field (email/phone/ID)
- `UnblockUser` - Unblock user by field (email/phone/ID)
- `GetUsers` - List users (paginated)
- `GetUserByField` - Get user by email/phone/ID

## Environment Variables

```bash
export GRPC_PORT=50051
export POSTGRES_HOST=localhost
export POSTGRES_PASSWORD=your_password
export AUTH_JWT_SECRET_KEY=your-secret-key
export REDIS_PASSWORD=your_redis_password
```

## Dependencies

Key packages:
- `gorm.io/gorm` - Database ORM
- `github.com/redis/go-redis/v9` - Redis client
- `google.golang.org/grpc` - gRPC framework
- `github.com/golang-jwt/jwt/v5` - JWT tokens
- `go.uber.org/zap` - Logging

That's it! Check the proto files in `/api/proto/auth/v1/` for detailed request/response schemas.