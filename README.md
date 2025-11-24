
# Qubool Kallyanam - Matrimonial Platform

A production-ready, cloud-native matrimonial platform built with Go microservices architecture, featuring matchmaking, real-time text chat, payment processing, and comprehensive admin controls.


## ✨ Key Features Implemented

### 1. **Microservices Architecture**
- **6 Independent Services**: Auth, User, Chat, Payment, Notification, Gateway
- **Clean Architecture**: Strict separation of layers (Domain, Use Case, Adapter, Handler)
- **gRPC Communication**: High-performance inter-service communication
- **Service Isolation**: Independent deployment and scaling capabilities


### 2. **Authentication & Authorization**
- **JWT-Based Authentication**: Secure access & refresh token management
- **Role-Based Access Control (RBAC)**: User, Premium User, and Admin roles
- **OTP Verification**: Email-based registration verification
- **Password Security**: Bcrypt hashing with salt
- **Session Management**: Redis-backed session store
- **Admin Controls**: User blocking/unblocking


### 3. **User Profile & Matchmaking System**
- **Comprehensive Profiles**: Demographics, education, profession, preferences
- **Partner Preferences**: Flexible search criteria (age, height, location, education, etc.)
- **Smart Recommendations**: Filter -based profile recommendations
- **Match Actions**: Like, Pass, and Mutual Match tracking
- **Photo Management**: Profile and additional photos with signed URLs
- **Privacy Controls**: Display order management for photos

### 4. **Real-Time Chat System**
- **One-to-One Conversations**: Private text messaging between users
- **WebSocket Support**: Real-time message delivery
- **Message Persistence**: MongoDB for conversation history
- **Premium Feature**: Chat restricted to premium users

### 5. **Payment & Subscription Management**
- **Razorpay Integration**: Secure payment gateway
- **Subscription Plans**: Admin-managed subscription tiers
- **Payment Verification**: Server-side signature verification
- **Payment History**: Complete transaction history for users
- **Auto-Activation**: Automatic premium status upon successful payment
- **Email Receipts**: Automated payment confirmation emails
- **Admin Dashboard**: View and manage completed payments


### 6. **Event based communication**
- **Asynchronous Communication**: Decoupled service interactions
- **Event Publishing**: Domain events for cross-service coordination
- **Event Subscriptions**: Multiple services consume relevant events
- **Dual Messaging Backend**: 
  - RabbitMQ for development
  - Google Cloud Pub/Sub for production
- **Event Types**: User lifecycle, profile updates, payments, match notifications

### 7. **Notification System**
- **Multi-Channel Notifications**: Email-based notifications
- **Event-Driven**: Automatically triggered by system events
- **Template Engine**: HTML email templates
- **Notification Types**:
  - OTP verification emails
  - Account deletion confirmations
  - Admin actions (blocking/unblocking notifications)
  - Match notifications (interest sent, mutual matches)
  - Payment successful notifications


### 8. **API Gateway Pattern**
- **Single Entry Point**: Unified REST API for all services
- **Protocol Translation**: HTTP/REST to gRPC conversion
- **Authentication Middleware**: JWT validation at gateway level
- **Request/Response Logging**: Structured logging with Zap
- **Error Normalization**: Consistent error responses across services

### 9. **Cloud-Native Infrastructure**
- **Kubernetes (GKE)**: Production deployments on Google Kubernetes Engine
- **Container Orchestration**: 
  - Deployment manifests for all services
  - Service discovery with Kubernetes Services
  - Ingress with Google-managed SSL certificates
  - Health checks and readiness probes
- **Workload Identity**: Secure GCP service authentication
- **Horizontal Scaling**: Auto-scaling based on load
- **Namespace Isolation**: Environment separation


### 10. **CI/CD Pipeline**
- **GitHub Actions Workflows**:
  - Automated Docker image builds
  - Google Artifact Registry integration
  - Selective service deployment
  - Database migration automation

### 11. **Database Strategy**
- **Polyglot Persistence**: Right database for each use case
  - **PostgreSQL**: User profiles, auth, payments, partner preferences, match data
  - **MongoDB**: Chat conversations and messages (document store)
  - **Redis**: Session management and OTP 
- **Migration Management**: Version-controlled SQL migrations


### 12. **Media Storage**
- **Google Cloud Storage**: Scalable object storage for photos
- **Signed URLs**: Temporary, secure upload/download URLs (15-min expiry)
- **Upload Confirmation Pattern**: Two-phase commit for uploads
- **Photo Deletion**: Lifecycle management of user photos
- **Development Emulator Support**: Local GCS emulator for testing

### 13. **Observability & Monitoring**
- **Structured Logging**: JSON-formatted logs with Zap
- **Prometheus Metrics**: Application and business metrics exposed
- **Request Tracing**: Request ID propagation across services
- **Health Checks**: Standard gRPC health checking protocol
- **Context Propagation**: User ID and request ID in gRPC metadata


### 14. **API Documentation**
- **Swagger/OpenAPI**: Interactive API documentation
- **Auto-Generated Docs**: Swagger UI at `/swagger`
- **Request/Response Examples**: Complete API contract documentation
- **Type Definitions**: Strongly-typed protobuf contracts


### 15. **Development Experience**
- **Docker Compose**: Complete local development environment
- **Configuration Management**: Viper-based config with env override
- **Hot Reload Ready**: Service-level configuration for development
- **Multi-Environment Support**: Development and production configs

## 🛠️ Technology Stack

### Backend Services
- **Go 1.23+**: High-performance, statically-typed language
- **gRPC**: Inter-service communication protocol
- **Protocol Buffers**: Service contract definitions

### Web Framework & API
- **Gin**: High-performance HTTP web framework for API Gateway
- **gorilla/websocket**: WebSocket implementation for real-time chat


### Databases
- **PostgreSQL**: Primary relational database
- **MongoDB**: Document store for chat data
- **Redis**: In-memory data store for OTP storage and sessions

### Message Brokers
- **RabbitMQ**: Development message broker
- **Google Cloud Pub/Sub**: Production message broker

### Cloud Services (Google Cloud Platform)
- **GKE (Google Kubernetes Engine)**: Container orchestration
- **Google Artifact Registry**: Container image storage
- **Google Cloud Storage**: Object storage for media
- **Workload Identity**: Secure service-to-service authentication
- **Cloud Pub/Sub**: Asynchronous messaging


### Security
- **JWT (golang-jwt/jwt)**: Token-based authentication
- **Bcrypt**: Password hashing
- **RBAC**: Role-based access control


### External Services
- **Razorpay**: Payment gateway integration
- **SMTP**: Email delivery for notifications

### DevOps & Infrastructure
- **Docker**: Containerization
- **Kubernetes**: Container orchestration
- **GitHub Actions**: CI/CD automation
- **Prometheus**: Metrics and monitoring

### Development Tools
- **GORM**: ORM for PostgreSQL
- **Viper**: Configuration management
- **Zap**: Structured logging
- **Swaggo**: Swagger documentation generation
- **golang-migrate**: Database migration tool
## 🚀 Getting Started

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 14+
- MongoDB 6+
- Redis 7+
- RabbitMQ 3.12+


## 🎯 API Endpoints

#### User Authentication
- `POST /auth/user/register` - User registration
- `POST /auth/user/verify` - user email verification using OTP
- `POST /auth/user/login` - User login
- `POST /auth/user/logout` - User logout (JWT + User role)
- `POST /auth/user/delete` - Delete account (JWT + User role)
- `POST /auth/user/refresh` - Refresh access token

#### Admin Authentication & User Management
- `POST /auth/admin/login` - Admin login
- `POST /auth/admin/logout` - Admin logout (JWT + Admin role)
- `POST /auth/admin/block-user` - Block user (JWT + Admin role)
- `POST /auth/admin/unblock-user` - Unblock user (JWT + Admin role)
- `GET /auth/admin/users` - List all users with pagination (JWT + Admin role)
- `GET /auth/admin/user` - Get user by email/phone/ID (JWT + Admin role)

### User Profile
- `GET /user/profile` - Get current user profile (JWT + User role)
- `PUT /user/profile` - Update entire user profile (JWT + User role)
- `PATCH /user/profile` - Partial profile update (JWT + User role)
- `GET /user/profiles/:profile_id` - Get full user details by profile ID (JWT + User role)
- `GET /user/profile-details/:profile_id` - Get full user details by profile ID (JWT + Admin role)


### Photo Management
- `POST /user/profile/profile-photo` - Get profile photo upload URL (JWT + User role)
- `POST /user/profile/profile-photo/confirm` - Confirm profile photo upload (JWT + User role)
- `DELETE /user/profile/profile-photo` - Delete profile photo (JWT + User role)
- `POST /user/profile/additional-photo` - Get additional photo upload URL (JWT + User role)
- `POST /user/profile/additional-photo/confirm` - Confirm additional photo upload (JWT + User role)
- `DELETE /user/profile/additional-photo/:display_order` - Delete additional photo (JWT + User role)
- `GET /user/profile/additional-photos` - List all additional photos (JWT + User role)

### Partner Preferences
- `POST /user/preference` - Create partner preferences (JWT + User role)
- `PATCH /user/preference` - Update partner preferences (JWT + User role)
- `GET /user/preference` - Get partner preferences (JWT + User role)

### Matchmaking
- `GET /user/recommendations` - Get profile recommendations (JWT + User role)
- `POST /user/match-action` - Record like/pass action (JWT + User role)
- `PUT /user/match-action` - Update match action (JWT + User role)
- `GET /user/matches/liked` - Get liked profiles (JWT + User role)
- `GET /user/matches/passed` - Get passed profiles (JWT + User role)
- `GET /user/matches/mutual` - Get mutual matches (JWT + User role)

### Chat (Premium Users Only)
- `POST /chat/conversation` - Create conversation (JWT + Premium User role)
- `GET /chat/conversation/:conversation_id/messages` - Get messages with pagination (JWT + Premium User role)
- `GET /chat/ws` - WebSocket connection for real-time messaging

### Payments (API Routes)
- `GET /payment/subscription-plans` - List active subscription plans (Public)
- `GET /payment/subscription-plan` - Get subscription plan by ID (Public)
- `POST /payment/subscription-plan` - Create subscription plan (JWT + Admin role)
- `PATCH /payment/subscription-plan` - Update subscription plan (JWT + Admin role)
- `POST /payment/order` - Create Razorpay payment order (JWT + User role)
- `GET /payment/subscriptions` - Get active subscription for user (JWT + User role)
- `GET /payment/payments-history` - Get payment history (JWT + User role)
- `GET /payment/admin/completed-payments` - Get all completed payments (JWT + Admin role)


### Payment (Web Pages - Public)
- `GET /payment/checkout` - Payment checkout page
- `GET /payment/verify` - Payment verification handler
- `GET /payment/success` - Payment success page
- `GET /payment/failed` - Payment failure page


> **Note**: All endpoints marked with "JWT + Role" require a valid JWT token in the Authorization header with the specified role. Premium User role is automatically granted upon successful subscription payment.

> For detailed request/response schemas and examples, visit the POSTMAN documentation at this [link](https://www.postman.com/manojsir/workspace/fawas-qk-project/collection/17830038-d9926cb2-542e-4677-a6f9-35cc0f60069d?action=share&creator=17830038)



## Author

- [@Mohamed Fawas](https://www.github.com/mohamedfawas)


