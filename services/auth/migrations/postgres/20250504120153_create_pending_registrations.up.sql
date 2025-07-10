CREATE TABLE IF NOT EXISTS pending_registrations (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    CHECK (expires_at > created_at)
);

-- Index to speed up cleanup of expired entries
CREATE INDEX IF NOT EXISTS idx_pending_reg_expires_at 
    ON pending_registrations (expires_at);
