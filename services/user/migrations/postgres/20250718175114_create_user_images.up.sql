CREATE TABLE IF NOT EXISTS user_images (
  id            BIGSERIAL PRIMARY KEY,
  user_id       UUID       NOT NULL,
  object_key    VARCHAR(500) NOT NULL,
  display_order SMALLINT   NOT NULL
      CONSTRAINT chk_display_order CHECK (display_order >= 1 AND display_order <= 3),
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
  deleted_at    TIMESTAMP WITH TIME ZONE
);

-- index to quickly lookup all active images for a user
CREATE INDEX IF NOT EXISTS idx_user_images_user_id
  ON user_images (user_id);

-- partial unique index: only counts rows where deleted_at IS NULL
CREATE UNIQUE INDEX IF NOT EXISTS uniq_active_img_order
  ON user_images (user_id, display_order)
  WHERE deleted_at IS NULL;
