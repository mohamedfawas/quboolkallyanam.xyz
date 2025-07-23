package dto

type UserProfilePutRequest struct {
	IsBride               bool    `json:"is_bride"`
	FullName              *string `json:"full_name"`
	DateOfBirth           *string `json:"date_of_birth"`
	HeightCm              *int    `json:"height_cm"`
	PhysicallyChallenged  bool    `json:"physically_challenged"`
	Community             *string `json:"community"`
	MaritalStatus         *string `json:"marital_status"`
	Profession            *string `json:"profession"`
	ProfessionType        *string `json:"profession_type"`
	HighestEducationLevel *string `json:"highest_education_level"`
	HomeDistrict          *string `json:"home_district"`
}

type UserProfilePatchRequest struct {
	IsBride               bool    `json:"is_bride,omitempty"`
	FullName              *string `json:"full_name,omitempty"`
	DateOfBirth           *string `json:"date_of_birth,omitempty"`
	HeightCm              *int    `json:"height_cm,omitempty"`
	PhysicallyChallenged  bool    `json:"physically_challenged,omitempty"`
	Community             *string `json:"community,omitempty"`
	MaritalStatus         *string `json:"marital_status,omitempty"`
	Profession            *string `json:"profession,omitempty"`
	ProfessionType        *string `json:"profession_type,omitempty"`
	HighestEducationLevel *string `json:"highest_education_level,omitempty"`
	HomeDistrict          *string `json:"home_district,omitempty"`
}

type PartnerPreferenceCreateRequest struct {
	MinAgeYears                *int     `json:"min_age_years"`
	MaxAgeYears                *int     `json:"max_age_years"`
	MinHeightCM                *int     `json:"min_height_cm"`
	MaxHeightCM                *int     `json:"max_height_cm"`
	AcceptPhysicallyChallenged bool     `json:"accept_physically_challenged"`
	PreferredCommunities       []string `json:"preferred_communities"`
	PreferredMaritalStatus     []string `json:"preferred_marital_status"`
	PreferredProfessions       []string `json:"preferred_professions"`
	PreferredProfessionTypes   []string `json:"preferred_profession_types"`
	PreferredEducationLevels   []string `json:"preferred_education_levels"`
	PreferredHomeDistricts     []string `json:"preferred_home_districts"`
}

type PartnerPreferencePatchRequest struct {
	MinAgeYears                *int      `json:"min_age_years,omitempty"`
	MaxAgeYears                *int      `json:"max_age_years,omitempty"`
	MinHeightCM                *int      `json:"min_height_cm,omitempty"`
	MaxHeightCM                *int      `json:"max_height_cm,omitempty"`
	AcceptPhysicallyChallenged *bool     `json:"accept_physically_challenged,omitempty"`
	PreferredCommunities       *[]string `json:"preferred_communities,omitempty"`
	PreferredMaritalStatus     *[]string `json:"preferred_marital_status,omitempty"`
	PreferredProfessions       *[]string `json:"preferred_professions,omitempty"`
	PreferredProfessionTypes   *[]string `json:"preferred_profession_types,omitempty"`
	PreferredEducationLevels   *[]string `json:"preferred_education_levels,omitempty"`
	PreferredHomeDistricts     *[]string `json:"preferred_home_districts,omitempty"`
}

type UpdateUserPartnerPreferencesResponse struct {
	Success bool `json:"success"`
}

type UpdateUserProfileResponse struct {
	Success bool `json:"success"`
}


type RecordMatchActionRequest struct {
	Action           string `json:"action"`
	TargetProfileID  uint   `json:"target_profile_id"`
}

type RecordMatchActionResponse struct {
	Success bool `json:"success"`
}