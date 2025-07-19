CREATE TABLE IF NOT EXISTS user_videos (
    id             BIGSERIAL    PRIMARY KEY,
    user_id               UUID           NOT NULL UNIQUE, -- one video per user id
    video_url             VARCHAR(500)   NOT NULL,
    object_key             VARCHAR(500)   NOT NULL,
    file_name             VARCHAR(255)   NOT NULL,
    file_size             BIGINT         NOT NULL,
    duration_seconds      INTEGER,
    created_at            TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_user_videos_user_id 
    FOREIGN KEY (user_id) REFERENCES user_profiles(user_id) ON DELETE CASCADE
);


CREATE INDEX idx_user_videos_user_id ON user_videos(user_id);