### Chat Service

Go-based chat microservice with gRPC API for creating conversations and sending/retrieving messages. Uses PostgreSQL for user projections, MongoDB for conversations/messages, and RabbitMQ (dev) or Google Pub/Sub (prod) for user-event ingestion.

### Features

- **Conversations**: Create one-to-one conversations
- **Messaging**: Send and list messages with pagination
- **User Projection**: Sync basic user data into PostgreSQL via events
- **Persistence**: MongoDB for conversations/messages, PostgreSQL for projections
- **Messaging**: RabbitMQ (development) or Pub/Sub (production)
- **gRPC**: Unary interceptor with structured error mapping
- **Logging**: Zap-based structured logs

### Quick Start

#### Prerequisites
- Go 1.23+
- PostgreSQL
- MongoDB
- RabbitMQ (development) or Google Cloud Pub/Sub (production)

#### Run Service
```bash
cd services/chat
go mod download
go run cmd/main.go
```

Service runs on port `50054` by default.

### Database Setup

- PostgreSQL (user projections)
```bash
# Run migration
psql -d your_db -f services/chat/migrations/postgres/20250719053841_create_user_projection.up.sql
```

- MongoDB (conversations/messages)
  - No migrations required. Ensure a database exists (default: `quboolKallyanam`).

### Configuration

Set environment variables or provide `config/config.yaml`. You can point to a custom config with `CONFIG_PATH`.

```yaml
environment: development

grpc:
  port: 50054

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: qubool_kallyanam_chat
  sslmode: disable
  timezone: UTC

mongodb:
  uri: "mongodb://localhost:27017"
  database: "quboolKallyanam"
  timeout: 10s

rabbitmq:
  dsn: "amqp://guest:guest@localhost:5672/"
  exchange_name: "qubool_kallyanam_events"

pubsub:
  project_id: "qubool-kallyanam-events"
```

### Environment Variables

```bash
export CONFIG_PATH=./config/config.yaml
export ENVIRONMENT=development
export GRPC_PORT=50054

export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=your_password
export POSTGRES_DBNAME=qubool_kallyanam_chat
export POSTGRES_SSLMODE=disable
export POSTGRES_TIMEZONE=UTC

export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="quboolKallyanam"

# Development messaging
export RABBITMQ_DSN="amqp://guest:guest@localhost:5672/"
export RABBITMQ_EXCHANGE_NAME="qubool_kallyanam_events"

# Production messaging (if ENVIRONMENT=production)
export PUBSUB_PROJECT_ID="qubool-kallyanam-events"
```

### gRPC API

Proto definitions: `api/proto/chat/v1/chat.proto`

- `CreateConversation(CreateConversationRequest) -> CreateConversationResponse`
- `SendMessage(SendMessageRequest) -> SendMessageResponse`
- `GetConversation(GetConversationRequest) -> GetConversationResponse`
- `GetMessagesByConversationId(GetMessagesByConversationIdRequest) -> GetMessagesByConversationIdResponse`

### Authentication / Metadata

All RPCs expect gRPC metadata:
- `request_id`: unique request ID
- `user_id`: authenticated user ID

Example (Go client):
```go
md := metadata.New(map[string]string{
  "request_id": "req-123",
  "user_id":    "user-456",
})
ctx := metadata.NewOutgoingContext(context.Background(), md)
resp, err := chatClient.SendMessage(ctx, &chatpbv1.SendMessageRequest{
  ConversationId: "...",
  Content: "Hello",
})
```

### Error Handling

- Application errors are returned as gRPC status with `google.rpc.ErrorInfo` details:
  - `reason`: app error code
  - `domain`: `chat`
  - `metadata.http_status_code`, `metadata.user_message`
- Unknown errors return `codes.Internal` with a generic message.

### Dependencies

- `google.golang.org/grpc` - gRPC framework
- `go.uber.org/zap` - Logging
- MongoDB and PostgreSQL clients via shared `pkg` modules