package dto

import "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"

///////////////////// USER PROFILE MANAGEMENT /////////////////////
type UserProfilePutRequest struct {
	IsBride               bool   `json:"is_bride"`
	FullName              string `json:"full_name"`
	DateOfBirth           string `json:"date_of_birth"`
	HeightCm              int    `json:"height_cm"`
	PhysicallyChallenged  bool   `json:"physically_challenged"`
	Community             string `json:"community"`
	MaritalStatus         string `json:"marital_status"`
	Profession            string `json:"profession"`
	ProfessionType        string `json:"profession_type"`
	HighestEducationLevel string `json:"highest_education_level"`
	HomeDistrict          string `json:"home_district"`
}

type UserProfilePatchRequest struct {
	IsBride               *bool   `json:"is_bride,omitempty"`
	FullName              *string `json:"full_name,omitempty"`
	DateOfBirth           *string `json:"date_of_birth,omitempty"`
	HeightCm              *int    `json:"height_cm,omitempty"`
	PhysicallyChallenged  *bool   `json:"physically_challenged,omitempty"`
	Community             *string `json:"community,omitempty"`
	MaritalStatus         *string `json:"marital_status,omitempty"`
	Profession            *string `json:"profession,omitempty"`
	ProfessionType        *string `json:"profession_type,omitempty"`
	HighestEducationLevel *string `json:"highest_education_level,omitempty"`
	HomeDistrict          *string `json:"home_district,omitempty"`
}

type UserProfileRecommendation struct {
	ID                int64  `json:"id"`
	FullName          string `json:"full_name"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Age               int    `json:"age"`
	HeightCm          int    `json:"height_cm"`
	MaritalStatus     string `json:"marital_status"`
	Profession        string `json:"profession"`
	HomeDistrict      string `json:"home_district"`
}

//////////////////// USER PROFILE PHOTO MANAGEMENT /////////////////////
type GetProfilePhotoUploadURLRequest struct {
	ContentType string `json:"content_type"`
}

type GetProfilePhotoUploadURLResponse struct {
	UploadURL        string `json:"upload_url"`
	ObjectKey        string `json:"object_key"`
	ExpiresInSeconds int32  `json:"expires_in_seconds"`
}

type ConfirmProfilePhotoUploadRequest struct {
	ObjectKey string `json:"object_key"`
}

type ConfirmProfilePhotoUploadResponse struct {
	Success bool `json:"success"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

type DeleteProfilePhotoRequest struct {
}

type DeleteProfilePhotoResponse struct {
	Success bool `json:"success"`
}	


/////////////////////////////// User Additional Photo Management ///////////////////////////////
type GetAdditionalPhotoUploadURLRequest struct {
	DisplayOrder int32 `json:"display_order"`
	ContentType string `json:"content_type"`
}

type GetAdditionalPhotoUploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	ExpiresInSeconds int32 `json:"expires_in_seconds"`
}


type ConfirmAdditionalPhotoUploadRequest struct {
	ObjectKey string `json:"object_key"`
}

type ConfirmAdditionalPhotoUploadResponse struct {
	Success bool `json:"success"`
	AdditionalPhotoURL string `json:"additional_photo_url"`
}

type DeleteAdditionalPhotoRequest struct {
	DisplayOrder int32 `json:"display_order"`
}

type DeleteAdditionalPhotoResponse struct {
	Success bool `json:"success"`
}

type GetAdditionalPhotosResponse struct {
	AdditionalPhotoURLs []string `json:"additional_photo_urls"`
}


/////////////////// PARTNER PREFERENCES MANAGEMENT /////////////////////
type PartnerPreferenceCreateRequest struct {
	MinAgeYears                int      `json:"min_age_years"`
	MaxAgeYears                int      `json:"max_age_years"`
	MinHeightCM                int      `json:"min_height_cm"`
	MaxHeightCM                int      `json:"max_height_cm"`
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

// similar to above given patch, but this is used in usecase => client layer , so no need to convert custom type to string type.
type UpdatePartnerPreferenceRequest struct {
	MinAgeYears                *int                         `json:"min_age_years,omitempty"`
	MaxAgeYears                *int                         `json:"max_age_years,omitempty"`
	MinHeightCM                *int                         `json:"min_height_cm,omitempty"`
	MaxHeightCM                *int                         `json:"max_height_cm,omitempty"`
	AcceptPhysicallyChallenged *bool                        `json:"accept_physically_challenged,omitempty"`
	PreferredCommunities       *[]validation.Community      `json:"preferred_communities,omitempty"`
	PreferredMaritalStatus     *[]validation.MaritalStatus  `json:"preferred_marital_status,omitempty"`
	PreferredProfessions       *[]validation.Profession     `json:"preferred_professions,omitempty"`
	PreferredProfessionTypes   *[]validation.ProfessionType `json:"preferred_profession_types,omitempty"`
	PreferredEducationLevels   *[]validation.EducationLevel `json:"preferred_education_levels,omitempty"`
	PreferredHomeDistricts     *[]validation.HomeDistrict   `json:"preferred_home_districts,omitempty"`
}

type UpdateUserPartnerPreferencesResponse struct {
	Success bool `json:"success"`
}

type UpdateUserProfileResponse struct {
	Success bool `json:"success"`
}

/////////////////// MATCH MAKING /////////////////////
type RecordMatchActionRequest struct {
	Action          string `json:"action"`
	TargetProfileID int64   `json:"target_profile_id"`
}

type RecordMatchActionResponse struct {
	Success bool `json:"success"`
}

type GetMatchRecommendationsRequest struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetMatchRecommendationsResponse struct {
	Profiles   []UserProfileRecommendation `json:"profiles"`
	Pagination PaginationInfo              `json:"pagination"`
}

type GetProfilesByMatchActionRequest struct {
	Action string `json:"action"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

type GetProfilesByMatchActionResponse struct {
	Profiles   []UserProfileRecommendation `json:"profiles"`
	Pagination PaginationInfo              `json:"pagination"`
}

/////////////////// PAGINATION /////////////////////
type PaginationInfo struct {
	TotalCount int64 `json:"total_count"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	HasMore    bool  `json:"has_more"`
}
