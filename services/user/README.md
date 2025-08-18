## User Service

Go-based microservice for user profiles, partner preferences, photos, and matchmaking over gRPC.

### Features

- **User profile**: Update/get profile
- **Photos**: Signed URLs for profile/additional photos; confirm/delete; list additional photos
- **Partner preferences**: Update/get preferences (flexible/strict)
- **Matchmaking**: Record actions (like/pass), recommendations, list by action
- **User details**: Fetch details by profile ID (admin-aware)
- **Events**: Consumes auth events for profile lifecycle

## Quick Start

### Prerequisites
- Go 1.23+
- PostgreSQL
- RabbitMQ (dev) or Google Pub/Sub (prod)
- Google Cloud Storage or a GCS emulator (dev via `MEDIA_STORAGE_ENDPOINT`)

### Run Service
```bash
cd services/user
go mod download
go run cmd/main.go
```

Service runs on port `50053` by default.

### Database Setup
Run migrations in order:
```bash
psql -d your_db -f migrations/postgres/20250718161605_create_user_profiles_and_partner_preferences.up.sql
psql -d your_db -f migrations/postgres/20250718175114_create_user_images.up.sql
psql -d your_db -f migrations/postgres/20250718175537_create_profile_matches.up.sql
psql -d your_db -f migrations/postgres/20250718175554_create_mutual_matches.up.sql
```

## Configuration

Set environment variables or create `internal/config/config.yaml`. You can override the config path using `CONFIG_PATH`.

Example `config.yaml`:
```yaml
environment: development

grpc:
  port: 50053

postgres:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: qubool_kallyanam_user
  sslmode: disable
  timezone: UTC

rabbitmq:
  dsn: amqp://guest:guest@localhost:5672/
  exchange_name: qubool_kallyanam_events

pubsub:
  project_id: qubool-kallyanam-events

media_storage:
  bucket: qubool-kallyanam-media
  credentials_file: secrets/gcs-service-account.json
  signer_email: ""
  private_key_path: secrets/signer_private_key.pem
  url_expiry: 15m
  endpoint: http://localhost:4443  # dev emulator; leave empty in prod
```

### Environment Variables
```bash
# General
export CONFIG_PATH=./internal/config/config.yaml
export ENVIRONMENT=development

# gRPC
export GRPC_PORT=50053

# PostgreSQL
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DBNAME=qubool_kallyanam_user
export POSTGRES_SSLMODE=disable
export POSTGRES_TIMEZONE=UTC

# Messaging (dev/prod)
export RABBITMQ_DSN=amqp://guest:guest@localhost:5672/
export RABBITMQ_EXCHANGE_NAME=qubool_kallyanam_events
export PUBSUB_PROJECT_ID=qubool-kallyanam-events

# Media storage
export MEDIA_STORAGE_BUCKET=qubool-kallyanam-media
export MEDIA_STORAGE_CREDENTIALS_FILE=secrets/gcs-service-account.json
export MEDIA_STORAGE_SIGNER_EMAIL=""
export MEDIA_STORAGE_PRIVATE_KEY_PATH=secrets/signer_private_key.pem
export MEDIA_STORAGE_URL_EXPIRY=15m
export MEDIA_STORAGE_ENDPOINT=http://localhost:4443
```

## gRPC API

See proto in `api/proto/user/v1/user.proto`. Service: `UserService`.

- `UpdateUserProfile` — Update current user profile
- `GetUserProfile` — Get current user profile
- `GetProfilePhotoUploadURL` — Signed URL for profile photo upload
- `ConfirmProfilePhotoUpload` — Confirm profile photo upload
- `DeleteProfilePhoto` — Remove profile photo
- `GetAdditionalPhotoUploadURL` — Signed URL for additional photo upload
- `ConfirmAdditionalPhotoUpload` — Confirm additional photo upload
- `DeleteAdditionalPhoto` — Remove additional photo by display order
- `GetAdditionalPhotos` — List additional photo URLs
- `UpdateUserPartnerPreferences` — Update partner preferences (supports granular/accept-all fields)
- `GetUserPartnerPreferences` — Get partner preferences
- `RecordMatchAction` — Record like/pass on a profile
- `GetMatchRecommendations` — Get recommended profiles (paginated)
- `GetProfilesByMatchAction` — List profiles by action (liked/passed)
- `GetUserDetailsByProfileID` — Fetch profile details by profile ID (admin-aware)

## Events

- Subscribes to auth events (creation, login, deletion) to initialize/update user profiles.
- Messaging backend: RabbitMQ in non-production; Google Pub/Sub in production.

## Dependencies

- `gorm.io/gorm` — PostgreSQL ORM
- `google.golang.org/grpc` — gRPC framework
- `go.uber.org/zap` — Logging
- Media storage via `pkg/mediastorage/gcs` (GCS or emulator)