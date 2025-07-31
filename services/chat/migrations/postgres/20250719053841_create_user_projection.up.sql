CREATE TABLE user_projection (
    user_uuid        UUID PRIMARY KEY, 
    user_profile_id  BIGINT UNIQUE NOT NULL,
    email            VARCHAR(255) NOT NULL,
    full_name        VARCHAR(255) NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_projection_user_profile_id ON user_projection(user_profile_id);