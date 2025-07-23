CREATE TABLE IF NOT EXISTS mutual_matches (
  id BIGSERIAL PRIMARY KEY,
  user_id_1 UUID NOT NULL,
  user_id_2 UUID NOT NULL,
  matched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  -- Status tracking
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMPTZ,

  -- Timestamps
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  -- Uniqueness constraint on sorted pair (before storing data make sure it's sorted)
  CONSTRAINT unique_mutual_match UNIQUE (user_id_1, user_id_2)
);

CREATE INDEX IF NOT EXISTS idx_mutual_matches_user1 ON mutual_matches (user_id_1) WHERE  is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_mutual_matches_user2 ON mutual_matches (user_id_2) WHERE  is_deleted = FALSE;

-- for clean up
CREATE INDEX IF NOT EXISTS idx_mutual_matches_is_deleted ON mutual_matches (is_deleted) WHERE is_deleted = TRUE;
