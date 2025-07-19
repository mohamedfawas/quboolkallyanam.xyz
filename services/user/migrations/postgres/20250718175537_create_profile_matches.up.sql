CREATE TABLE IF NOT EXISTS profile_matches (
  id BIGSERIAL PRIMARY KEY,
  user_id UUID NOT NULL,
  target_id UUID NOT NULL,
  is_liked BOOLEAN NOT NULL DEFAULT false, -- only two options 'liked' or 'passed'
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ,
  CONSTRAINT unique_profile_match UNIQUE(user_id, target_id)
);

-- User can know whom they have liked or passed
CREATE INDEX IF NOT EXISTS idx_profile_matches_user_liked ON profile_matches(user_id, is_liked) WHERE is_deleted = false;
-- User can know whom has liked them
CREATE INDEX IF NOT EXISTS idx_profile_matches_target_liked ON profile_matches(target_id, is_liked) WHERE is_deleted = false;
-- for clean up
CREATE INDEX IF NOT EXISTS idx_profile_matches_is_deleted ON profile_matches(is_deleted) WHERE is_deleted = true;
