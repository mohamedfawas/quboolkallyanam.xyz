CREATE TABLE IF NOT EXISTS user_images (
  id                    BIGSERIAL      PRIMARY KEY,
  user_id        UUID         NOT NULL REFERENCES user_profiles(user_id) ON DELETE CASCADE,
  object_key     VARCHAR(255) NOT NULL,
  display_order  SMALLINT     NOT NULL CHECK (display_order BETWEEN 1 AND 3),
  created_at            TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

  CONSTRAINT unique_user_images_order 
    UNIQUE (user_id, display_order)
);


CREATE INDEX IF NOT EXISTS idx_user_images_user_id ON user_images(user_id);
