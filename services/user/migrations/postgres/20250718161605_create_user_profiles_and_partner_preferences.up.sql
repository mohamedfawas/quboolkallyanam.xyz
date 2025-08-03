CREATE TABLE IF NOT EXISTS user_profiles (
  id                      BIGSERIAL      PRIMARY KEY,    -- local PK
  user_id                 UUID           NOT NULL UNIQUE,       -- from Auth service
  is_bride                BOOLEAN        NOT NULL DEFAULT FALSE,
  full_name               VARCHAR(200),
  email                   VARCHAR(255),
  phone                   VARCHAR(20),
  date_of_birth           DATE,
  height_cm               SMALLINT NOT NULL CHECK (height_cm > 0),
  physically_challenged   BOOLEAN NOT NULL DEFAULT FALSE,
  community               VARCHAR(255),
  marital_status          VARCHAR(255),
  profession              VARCHAR(255),
  profession_type         VARCHAR(255),
  highest_education_level VARCHAR(255),
  home_district           VARCHAR(255),
  profile_image_key     VARCHAR(255),
  last_login              TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  created_at              TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  updated_at              TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  deleted_at              TIMESTAMPTZ  
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_last_login ON user_profiles (last_login DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_profiles_community  ON user_profiles(community)      WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_profiles_marital_status ON user_profiles(marital_status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_profiles_dob            ON user_profiles(date_of_birth)  WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_profiles_height         ON user_profiles(height_cm)      WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_profiles_profile_image_key ON user_profiles(profile_image_key);

-- Partner preferences table
CREATE TABLE IF NOT EXISTS partner_preferences (
	id                        BIGSERIAL       PRIMARY KEY,
	user_profile_id           BIGINT          NOT NULL REFERENCES user_profiles(id),
	
	-- Range preferences
	min_age_years             SMALLINT            NOT NULL DEFAULT 18,
	max_age_years             SMALLINT            NOT NULL DEFAULT 100,
	min_height_cm             SMALLINT            NOT NULL DEFAULT 130,
	max_height_cm             SMALLINT            NOT NULL DEFAULT 220,

	-- Simple boolean preferences
	accept_physically_challenged BOOLEAN      NOT NULL DEFAULT TRUE,
	
  preferred_communities        TEXT[]        NOT NULL DEFAULT '{}',
	preferred_marital_status     TEXT[]        NOT NULL DEFAULT '{}',
	preferred_professions        TEXT[]        NOT NULL DEFAULT '{}',
	preferred_profession_types   TEXT[]        NOT NULL DEFAULT '{}',
	preferred_education_levels   TEXT[]        NOT NULL DEFAULT '{}',
	preferred_home_districts     TEXT[]        NOT NULL DEFAULT '{}',

  
	created_at                TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
  updated_at                TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
  	
  -- Constraints
	CONSTRAINT age_range_valid CHECK (max_age_years >= min_age_years),
	CONSTRAINT height_range_valid CHECK (max_height_cm >= min_height_cm),
	CONSTRAINT user_preferences_unique UNIQUE (user_profile_id)
);


-- PRIMARY MATCHING INDEX 
CREATE INDEX IF NOT EXISTS idx_user_profiles_matching 
ON user_profiles (deleted_at, is_bride, last_login DESC) 
WHERE deleted_at IS NULL;

--  AGE RANGE INDEX (Critical for age filtering)
CREATE INDEX IF NOT EXISTS idx_user_profiles_age_range 
ON user_profiles (deleted_at, date_of_birth, last_login DESC) 
WHERE deleted_at IS NULL;

--  HEIGHT RANGE INDEX (Critical for height filtering)
CREATE INDEX IF NOT EXISTS idx_user_profiles_height_range 
ON user_profiles (deleted_at, height_cm, last_login DESC) 
WHERE deleted_at IS NULL;

-- PHYSICALLY CHALLENGED INDEX (For boolean filtering)
CREATE INDEX IF NOT EXISTS idx_user_profiles_physically_challenged 
ON user_profiles (deleted_at, physically_challenged, last_login DESC) 
WHERE deleted_at IS NULL;


-- B-tree indexes on boolean flags
CREATE INDEX IF NOT EXISTS idx_partner_pref_accept_all_communities
  ON partner_preferences (accept_all_communities);
CREATE INDEX IF NOT EXISTS idx_partner_pref_accept_all_marital_status
  ON partner_preferences (accept_all_marital_status);
CREATE INDEX IF NOT EXISTS idx_partner_pref_accept_all_professions
  ON partner_preferences (accept_all_professions);
CREATE INDEX IF NOT EXISTS idx_partner_pref_accept_all_profession_types
  ON partner_preferences (accept_all_profession_types);
CREATE INDEX IF NOT EXISTS idx_partner_pref_accept_all_education_levels
  ON partner_preferences (accept_all_education_levels);
CREATE INDEX IF NOT EXISTS idx_partner_pref_accept_all_home_districts
  ON partner_preferences (accept_all_home_districts);

-- GIN indexes on array columns
CREATE INDEX IF NOT EXISTS idx_partner_pref_communities_gin
  ON partner_preferences USING gin (preferred_communities);
CREATE INDEX IF NOT EXISTS idx_partner_pref_marital_status_gin
  ON partner_preferences USING gin (preferred_marital_status);
CREATE INDEX IF NOT EXISTS idx_partner_pref_professions_gin
  ON partner_preferences USING gin (preferred_professions);
CREATE INDEX IF NOT EXISTS idx_partner_pref_profession_types_gin
  ON partner_preferences USING gin (preferred_profession_types);
CREATE INDEX IF NOT EXISTS idx_partner_pref_education_levels_gin
  ON partner_preferences USING gin (preferred_education_levels);
CREATE INDEX IF NOT EXISTS idx_partner_pref_home_districts_gin
  ON partner_preferences USING gin (preferred_home_districts);