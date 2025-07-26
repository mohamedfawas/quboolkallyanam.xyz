-- 1. ENUM types
CREATE TYPE community_enum AS ENUM (
  'sunni','mujahid','tabligh','jamate_islami','shia','muslim','not_mentioned'
);

CREATE TYPE marital_status_enum AS ENUM (
  'never_married','divorced','nikkah_divorce','widowed','not_mentioned'
);

CREATE TYPE profession_enum AS ENUM (
  'student','doctor','engineer','farmer','teacher','not_mentioned'
);

CREATE TYPE profession_type_enum AS ENUM (
  'full_time','part_time','freelance','self_employed','not_working','not_mentioned'
);

CREATE TYPE education_level_enum AS ENUM (
  'less_than_high_school','high_school','higher_secondary',
  'under_graduation','post_graduation','not_mentioned'
);

CREATE TYPE home_district_enum AS ENUM (
  'thiruvananthapuram','kollam','pathanamthitta','alappuzha',
  'kottayam','ernakulam','thrissur','palakkad',
  'malappuram','kozhikode','wayanad','kannur',
  'kasaragod','idukki','not_mentioned'
);



-- 2. user_profiles table

CREATE TABLE IF NOT EXISTS user_profiles (
  id                      BIGSERIAL      PRIMARY KEY,    -- local PK
  user_id                 UUID           NOT NULL UNIQUE,       -- from Auth service
  is_bride                BOOLEAN        NOT NULL DEFAULT FALSE,
  full_name               VARCHAR(200)   NULL,
  email                   VARCHAR(255)   ,
  phone                   VARCHAR(20)    ,
  date_of_birth           DATE           NULL,
  height_cm               INT            CHECK (height_cm BETWEEN 130 AND 220),
  physically_challenged    BOOLEAN        NOT NULL DEFAULT FALSE,
  community               community_enum NULL,
  marital_status          marital_status_enum NULL,
  profession              profession_enum NULL,
  profession_type         profession_type_enum NULL,
  highest_education_level education_level_enum NULL,
  home_district           home_district_enum NULL,
  profile_picture_url     VARCHAR(255) NULL,
  last_login              TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  created_at              TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  updated_at              TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  is_deleted              BOOLEAN        NOT NULL DEFAULT FALSE,
  deleted_at              TIMESTAMPTZ
);


CREATE INDEX IF NOT EXISTS idx_user_profiles_last_login ON user_profiles (last_login DESC) WHERE is_deleted = false;
CREATE INDEX IF NOT EXISTS idx_user_profiles_community  ON user_profiles(community)      WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_profiles_marital_status ON user_profiles(marital_status) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_profiles_dob            ON user_profiles(date_of_birth)  WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_profiles_height         ON user_profiles(height_cm)      WHERE is_deleted = FALSE;



-- Partner preferences table
CREATE TABLE IF NOT EXISTS partner_preferences (
	id                        BIGSERIAL       PRIMARY KEY,
	user_profile_id           BIGINT          NOT NULL REFERENCES user_profiles(id),
	
	-- Range preferences
	min_age_years             INT             CHECK (min_age_years BETWEEN 18 AND 100),
	max_age_years             INT             CHECK (max_age_years BETWEEN 18 AND 100),
	min_height_cm             INT             CHECK (min_height_cm BETWEEN 130 AND 220),
	max_height_cm             INT             CHECK (max_height_cm BETWEEN 130 AND 220),
	  
	-- Simple boolean preferences
	accept_physically_challenged BOOLEAN      NOT NULL DEFAULT TRUE,
	
	-- JSONB preferences 
    preferred_communities     JSONB           NOT NULL DEFAULT '[]'::jsonb,
    preferred_marital_status  JSONB           NOT NULL DEFAULT '[]'::jsonb,
    preferred_professions     JSONB           NOT NULL DEFAULT '[]'::jsonb,
    preferred_profession_types JSONB          NOT NULL DEFAULT '[]'::jsonb,
    preferred_education_levels JSONB          NOT NULL DEFAULT '[]'::jsonb,
    preferred_home_districts  JSONB           NOT NULL DEFAULT '[]'::jsonb,
	
	created_at                TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
  	updated_at                TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
  	
  	-- Constraints
	CONSTRAINT age_range_valid CHECK (max_age_years >= min_age_years),
	CONSTRAINT height_range_valid CHECK (max_height_cm >= min_height_cm),
	CONSTRAINT user_preferences_unique UNIQUE (user_profile_id)
);


-- PRIMARY MATCHING INDEX 
CREATE INDEX IF NOT EXISTS idx_user_profiles_matching 
ON user_profiles (is_deleted, is_bride, last_login DESC) 
WHERE is_deleted = false;

--  AGE RANGE INDEX (Critical for age filtering)
CREATE INDEX IF NOT EXISTS idx_user_profiles_age_range 
ON user_profiles (is_deleted, date_of_birth, last_login DESC) 
WHERE is_deleted = false;

--  HEIGHT RANGE INDEX (Critical for height filtering)
CREATE INDEX IF NOT EXISTS idx_user_profiles_height_range 
ON user_profiles (is_deleted, height_cm, last_login DESC) 
WHERE is_deleted = false;

-- PHYSICALLY CHALLENGED INDEX (For boolean filtering)
CREATE INDEX IF NOT EXISTS idx_user_profiles_physically_challenged 
ON user_profiles (is_deleted, physically_challenged, last_login DESC) 
WHERE is_deleted = false;


-- GIN indexes for JSONB fields (essential for performance)
CREATE INDEX IF NOT EXISTS idx_partner_pref_communities 
ON partner_preferences USING GIN (preferred_communities);

CREATE INDEX IF NOT EXISTS idx_partner_pref_marital_status 
ON partner_preferences USING GIN (preferred_marital_status);

CREATE INDEX IF NOT EXISTS idx_partner_pref_professions 
ON partner_preferences USING GIN (preferred_professions);

CREATE INDEX IF NOT EXISTS idx_partner_pref_education 
ON partner_preferences USING GIN (preferred_education_levels);

CREATE INDEX IF NOT EXISTS idx_partner_pref_districts 
ON partner_preferences USING GIN (preferred_home_districts);
