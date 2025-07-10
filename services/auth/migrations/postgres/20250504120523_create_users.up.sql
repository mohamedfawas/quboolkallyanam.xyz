
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL ,
    phone VARCHAR(20) NOT NULL ,
    password_hash VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    premium_until TIMESTAMP,
    last_login_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL ,
    updated_at TIMESTAMP NOT NULL ,
    deleted_at TIMESTAMP
);

-- Create partial unique constraints that only apply to non-deleted records
CREATE UNIQUE INDEX users_email_unique_active ON users (email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX users_phone_unique_active ON users (phone) WHERE deleted_at IS NULL;

-- Create indexes for blocked users
CREATE INDEX users_blocked_email_idx ON users (email) WHERE is_blocked = true;
CREATE INDEX users_blocked_phone_idx ON users (phone) WHERE is_blocked = true;