### Payment Service

Go-based payment microservice with gRPC API for subscription plans, order creation, Razorpay payment verification, and subscription management.

### Features

- **Payment Orders**: Create Razorpay orders for subscription plans
- **Payment Page Data**: Retrieve front-end payment initialization data
- **Payment Verification**: Verify Razorpay signatures, activate subscriptions, and publish events
- **Subscriptions**:
  - Admin create/update subscription plans
  - Get active plans
  - Get user’s active subscription
- **History & Admin Views**:
  - User payment history
  - Completed payment details (paginated)
- **Messaging**: RabbitMQ (dev) or Google Cloud Pub/Sub (prod)
- **Config**: Viper-based, env-first with optional YAML

### Quick Start

#### Prerequisites
- Go 1.23+
- PostgreSQL
- RabbitMQ (dev) or Google Cloud Pub/Sub (prod)
- Razorpay account and API keys
- Optional: SMTP for email notifications

#### Run Service
```bash
cd services/payment
go mod download
go run cmd/main.go
```
Service runs on port `50055` by default.

### Database Setup

Run migrations in order:
```bash
psql -d your_db -f migrations/postgres/20250711142519_create_subscription_plans.up.sql
psql -d your_db -f migrations/postgres/20250711143742_create_subscriptions.up.sql
psql -d your_db -f migrations/postgres/20250711150246_create_payments.up.sql
```

### Configuration

Set environment variables or create `config/config.yaml` (the service auto-loads from `CONFIG_PATH` or `./config/config.yaml`).

```yaml
environment: development

grpc:
  port: 50055

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: qubool_kallyanam_payment
  sslmode: disable
  timezone: UTC

razorpay:
  key_id: "your-key-id"
  key_secret: "your-key-secret"

email:
  smtp_host: smtp.gmail.com
  smtp_port: 587
  smtp_username: ""
  smtp_password: ""
  from_email: "noreply@qubool-kallyanam.xyz"
  from_name: "Qubool Kallyanam"

rabbitmq:
  dsn: "amqp://guest:guest@localhost:5672/"
  exchange_name: "qubool_kallyanam_events"

pubsub:
  project_id: "qubool-kallyanam-events"
```

### API Endpoints

Proto: `api/proto/payment/v1/payment.proto`

- `CreatePaymentOrder` — Create Razorpay order for a plan
- `ShowPaymentPage` — Fetch front-end payment initialization data
- `VerifyPayment` — Verify Razorpay payment and activate subscription
- `CreateOrUpdateSubscriptionPlan` — Admin upsert subscription plan
- `GetSubscriptionPlan` — Retrieve a plan by ID
- `GetActiveSubscriptionPlans` — List active plans
- `GetActiveSubscriptionByUserID` — Current user’s active subscription
- `GetPaymentHistory` — Current user’s payment history
- `GetCompletedPaymentDetails` — Admin paginated completed payments

### Environment Variables

```bash
export ENVIRONMENT=development
export CONFIG_PATH=./config/config.yaml

export GRPC_PORT=50055

export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DBNAME=qubool_kallyanam_payment
export POSTGRES_SSLMODE=disable
export POSTGRES_TIMEZONE=UTC

export RAZORPAY_KEY_ID=your-key-id
export RAZORPAY_KEY_SECRET=your-key-secret

export EMAIL_SMTP_HOST=smtp.gmail.com
export EMAIL_SMTP_PORT=587
export EMAIL_SMTP_USERNAME=""
export EMAIL_SMTP_PASSWORD=""
export EMAIL_FROM_EMAIL=noreply@qubool-kallyanam.xyz
export EMAIL_FROM_NAME="Qubool Kallyanam"

export RABBITMQ_DSN="amqp://guest:guest@localhost:5672/"
export RABBITMQ_EXCHANGE_NAME=qubool_kallyanam_events

export PUBSUB_PROJECT_ID=qubool-kallyanam-events
```

### Dependencies

Key packages:
- `gorm.io/gorm` — Database ORM
- `google.golang.org/grpc` — gRPC framework
- `go.uber.org/zap` — Structured logging
- `github.com/spf13/viper` — Configuration
- `github.com/razorpay/razorpay-go` — Razorpay API
- RabbitMQ and Google Pub/Sub clients via internal adapters

That’s it! Check `api/proto/payment/v1/payment.proto` for detailed request/response schemas.