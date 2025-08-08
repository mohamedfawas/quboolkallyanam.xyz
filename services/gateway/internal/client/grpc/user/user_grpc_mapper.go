package user

import (
	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/pointerutil"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

////////////////////////////// Update User Profile //////////////////////////////

func MapUpdateUserProfileRequest(req dto.UserProfilePatchRequest) *userpbv1.UpdateUserProfileRequest {
	grpcReq := &userpbv1.UpdateUserProfileRequest{}

	grpcReq.IsBride = &wrapperspb.BoolValue{Value: pointerutil.GetBoolValue(req.IsBride)}
	grpcReq.PhysicallyChallenged = &wrapperspb.BoolValue{Value: pointerutil.GetBoolValue(req.PhysicallyChallenged)}

	if req.FullName != nil {
		grpcReq.FullName = &wrapperspb.StringValue{Value: pointerutil.GetStringValue(req.FullName)}
	}

	if req.DateOfBirth != nil {
		grpcReq.DateOfBirth = &wrapperspb.StringValue{Value: pointerutil.GetStringValue(req.DateOfBirth)}
	}
	if req.Community != nil {
		grpcReq.Community = &wrapperspb.StringValue{Value: pointerutil.GetStringValue(req.Community)}
	}
	if req.MaritalStatus != nil {
		grpcReq.MaritalStatus = &wrapperspb.StringValue{Value: pointerutil.GetStringValue(req.MaritalStatus)}
	}
	if req.Profession != nil {
		grpcReq.Profession = &wrapperspb.StringValue{Value: pointerutil.GetStringValue(req.Profession)}
	}
	if req.ProfessionType != nil {
		grpcReq.ProfessionType = &wrapperspb.StringValue{Value: pointerutil.GetStringValue(req.ProfessionType)}
	}
	if req.HighestEducationLevel != nil {
		grpcReq.HighestEducationLevel = &wrapperspb.StringValue{Value: pointerutil.GetStringValue(req.HighestEducationLevel)}
	}
	if req.HomeDistrict != nil {
		grpcReq.HomeDistrict = &wrapperspb.StringValue{Value: pointerutil.GetStringValue(req.HomeDistrict)}
	}

	if req.HeightCm != nil {
		grpcReq.HeightCm = &wrapperspb.Int32Value{Value: int32(pointerutil.GetIntValue(req.HeightCm))}
	}

	return grpcReq
}

func MapUpdateUserProfileResponse(resp *userpbv1.UpdateUserProfileResponse) *dto.UpdateUserProfileResponse {
	return &dto.UpdateUserProfileResponse{
		Success: resp.Success.GetValue(),
	}
}

////////////////////////////// Get User Profile //////////////////////////////

func MapGetUserProfileResponse(resp *userpbv1.GetUserProfileResponse) *dto.UserProfileRecommendation {
	p := resp.GetProfile()
	if p == nil {
		return nil
	}
	return &dto.UserProfileRecommendation{
		ID:                p.Id,
		FullName:          p.FullName,
		ProfilePictureURL: p.ProfilePictureUrl,
		Age:               int(p.Age),
		HeightCm:          int(p.HeightCm),
		MaritalStatus:     p.MaritalStatus,
		Profession:        p.Profession,
		HomeDistrict:      p.HomeDistrict,
	}
}

////////////////////////////// Get Profile Photo Upload URL //////////////////////////////

func MapGetProfilePhotoUploadURLRequest(req dto.GetProfilePhotoUploadURLRequest) *userpbv1.GetProfilePhotoUploadURLRequest {
	return &userpbv1.GetProfilePhotoUploadURLRequest{
		ContentType: &wrapperspb.StringValue{Value: req.ContentType},
	}
}

func MapGetProfilePhotoUploadURLResponse(resp *userpbv1.GetProfilePhotoUploadURLResponse) *dto.GetProfilePhotoUploadURLResponse {
	return &dto.GetProfilePhotoUploadURLResponse{
		UploadURL: resp.UploadUrl.GetValue(),
		ObjectKey: resp.ObjectKey.GetValue(),
		ExpiresInSeconds: int32(resp.ExpiresInSeconds.GetValue()),
	}
}

////////////////////////////// Confirm Profile Photo Upload //////////////////////////////

func MapConfirmProfilePhotoUploadRequest(req dto.ConfirmProfilePhotoUploadRequest) *userpbv1.ConfirmProfilePhotoUploadRequest {
	return &userpbv1.ConfirmProfilePhotoUploadRequest{
		ObjectKey: &wrapperspb.StringValue{Value: req.ObjectKey},
		FileSize: &wrapperspb.Int64Value{Value: req.FileSize},
	}
}

func MapConfirmProfilePhotoUploadResponse(resp *userpbv1.ConfirmProfilePhotoUploadResponse) *dto.ConfirmProfilePhotoUploadResponse {
	return &dto.ConfirmProfilePhotoUploadResponse{
		Success: resp.Success.GetValue(),
		ProfilePictureURL: resp.ProfilePictureUrl.GetValue(),
	}
}

////////////////////////////// Delete Profile Photo //////////////////////////////

func MapDeleteProfilePhotoRequest(req dto.DeleteProfilePhotoRequest) *userpbv1.DeleteProfilePhotoRequest {
	return &userpbv1.DeleteProfilePhotoRequest{}
}

func MapDeleteProfilePhotoResponse(resp *userpbv1.DeleteProfilePhotoResponse) *dto.DeleteProfilePhotoResponse {
	return &dto.DeleteProfilePhotoResponse{
		Success: resp.Success.GetValue(),
	}
}

////////////////////////////// Get Additional Photo Upload URL //////////////////////////////

func MapGetAdditionalPhotoUploadURLRequest(req dto.GetAdditionalPhotoUploadURLRequest) *userpbv1.GetAdditionalPhotoUploadURLRequest {
	return &userpbv1.GetAdditionalPhotoUploadURLRequest{
		DisplayOrder: &wrapperspb.Int32Value{Value: int32(req.DisplayOrder)},
		ContentType: &wrapperspb.StringValue{Value: req.ContentType},
	}
}

func MapGetAdditionalPhotoUploadURLResponse(resp *userpbv1.GetAdditionalPhotoUploadURLResponse) *dto.GetAdditionalPhotoUploadURLResponse {
	return &dto.GetAdditionalPhotoUploadURLResponse{
		UploadURL: resp.UploadUrl.GetValue(),
		ObjectKey: resp.ObjectKey.GetValue(),
		ExpiresInSeconds: int32(resp.ExpiresInSeconds.GetValue()),
	}
}

////////////////////////////// Confirm Additional Photo Upload //////////////////////////////

func MapConfirmAdditionalPhotoUploadRequest(req dto.ConfirmAdditionalPhotoUploadRequest) *userpbv1.ConfirmAdditionalPhotoUploadRequest {
	return &userpbv1.ConfirmAdditionalPhotoUploadRequest{
		ObjectKey: &wrapperspb.StringValue{Value: req.ObjectKey},
		FileSize: &wrapperspb.Int64Value{Value: req.FileSize},
	}
}

func MapConfirmAdditionalPhotoUploadResponse(resp *userpbv1.ConfirmAdditionalPhotoUploadResponse) *dto.ConfirmAdditionalPhotoUploadResponse {
	return &dto.ConfirmAdditionalPhotoUploadResponse{
		Success: resp.Success.GetValue(),
		AdditionalPhotoURL: resp.AdditionalPhotoUrl.GetValue(),
	}
}


////////////////////////////// Delete Additional Photo //////////////////////////////

func MapDeleteAdditionalPhotoRequest(req dto.DeleteAdditionalPhotoRequest) *userpbv1.DeleteAdditionalPhotoRequest {
	return &userpbv1.DeleteAdditionalPhotoRequest{
		DisplayOrder: &wrapperspb.Int32Value{Value: int32(req.DisplayOrder)},
	}
}


func MapDeleteAdditionalPhotoResponse(resp *userpbv1.DeleteAdditionalPhotoResponse) *dto.DeleteAdditionalPhotoResponse {
	return &dto.DeleteAdditionalPhotoResponse{
		Success: resp.Success.GetValue(),
	}
}

////////////////////////////// Update User Partner Preferences //////////////////////////////

func MapUpdateUserPartnerPreferencesRequest(operationType string, req dto.UpdatePartnerPreferenceRequest) *userpbv1.UpdateUserPartnerPreferencesRequest {
	grpcReq := &userpbv1.UpdateUserPartnerPreferencesRequest{
		OperationType: &wrapperspb.StringValue{Value: operationType},
	}

	if req.MinAgeYears != nil {
		grpcReq.MinAgeYears = &wrapperspb.Int32Value{Value: int32(*req.MinAgeYears)}
	}

	if req.MaxAgeYears != nil {
		grpcReq.MaxAgeYears = &wrapperspb.Int32Value{Value: int32(*req.MaxAgeYears)}
	}

	if req.MinHeightCM != nil {
		grpcReq.MinHeightCm = &wrapperspb.Int32Value{Value: int32(*req.MinHeightCM)}
	}

	if req.MaxHeightCM != nil {
		grpcReq.MaxHeightCm = &wrapperspb.Int32Value{Value: int32(*req.MaxHeightCM)}
	}

	if req.AcceptPhysicallyChallenged != nil {
		grpcReq.AcceptPhysicallyChallenged = &wrapperspb.BoolValue{Value: *req.AcceptPhysicallyChallenged}
	}

	if req.PreferredCommunities != nil {
		communities := make([]string, len(*req.PreferredCommunities))
		for i, community := range *req.PreferredCommunities {
			communities[i] = string(community)
		}
		grpcReq.PreferredCommunities = communities
	}

	if req.PreferredMaritalStatus != nil {
		maritalStatuses := make([]string, len(*req.PreferredMaritalStatus))
		for i, status := range *req.PreferredMaritalStatus {
			maritalStatuses[i] = string(status)
		}
		grpcReq.PreferredMaritalStatus = maritalStatuses
	}

	if req.PreferredProfessions != nil {
		professions := make([]string, len(*req.PreferredProfessions))
		for i, profession := range *req.PreferredProfessions {
			professions[i] = string(profession)
		}
		grpcReq.PreferredProfessions = professions
	}

	if req.PreferredProfessionTypes != nil {
		professionTypes := make([]string, len(*req.PreferredProfessionTypes))
		for i, professionType := range *req.PreferredProfessionTypes {
			professionTypes[i] = string(professionType)
		}
		grpcReq.PreferredProfessionTypes = professionTypes
	}

	if req.PreferredEducationLevels != nil {
		educationLevels := make([]string, len(*req.PreferredEducationLevels))
		for i, level := range *req.PreferredEducationLevels {
			educationLevels[i] = string(level)
		}
		grpcReq.PreferredEducationLevels = educationLevels
	}

	if req.PreferredHomeDistricts != nil {
		homeDistricts := make([]string, len(*req.PreferredHomeDistricts))
		for i, district := range *req.PreferredHomeDistricts {
			homeDistricts[i] = string(district)
		}
		grpcReq.PreferredHomeDistricts = homeDistricts
	}

	return grpcReq
}

func MapUpdateUserPartnerPreferencesResponse(resp *userpbv1.UpdateUserPartnerPreferencesResponse) *dto.UpdateUserPartnerPreferencesResponse {
	return &dto.UpdateUserPartnerPreferencesResponse{
		Success: resp.Success.GetValue(),
	}
}


////////////////////////////// Record Match Action //////////////////////////////

func MapRecordMatchActionRequest(req dto.RecordMatchActionRequest) *userpbv1.RecordMatchActionRequest {
	return &userpbv1.RecordMatchActionRequest{
		Action:          req.Action,
		TargetProfileId: req.TargetProfileID,
	}
}

func MapRecordMatchActionResponse(resp *userpbv1.RecordMatchActionResponse) *dto.RecordMatchActionResponse {
	return &dto.RecordMatchActionResponse{
		Success: resp.Success,
	}
}



////////////////////////////// Get Match Recommendations //////////////////////////////

func MapGetMatchRecommendationsRequest(req dto.GetMatchRecommendationsRequest) *userpbv1.GetMatchRecommendationsRequest {
	return &userpbv1.GetMatchRecommendationsRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}

func MapGetMatchRecommendationsResponse(resp *userpbv1.GetMatchRecommendationsResponse) *dto.GetMatchRecommendationsResponse {
	profiles := make([]dto.UserProfileRecommendation, len(resp.Profiles))
	for i, profile := range resp.Profiles {
		profiles[i] = dto.UserProfileRecommendation{
			ID:                profile.Id,
			FullName:          profile.FullName,
			ProfilePictureURL: profile.ProfilePictureUrl,
			Age:               int(profile.Age),
			HeightCm:          int(profile.HeightCm),
			MaritalStatus:     profile.MaritalStatus,
			Profession:        profile.Profession,
			HomeDistrict:      profile.HomeDistrict,
		}
	}

	pagination := dto.PaginationInfo{
		TotalCount: resp.Pagination.TotalCount,
		Limit:      int(resp.Pagination.Limit),
		Offset:     int(resp.Pagination.Offset),
		HasMore:    resp.Pagination.HasMore,
	}

	return &dto.GetMatchRecommendationsResponse{
		Profiles:   profiles,
		Pagination: pagination,
	}
}


////////////////////////////// Get Profiles By Match Action //////////////////////////////

func MapGetProfilesByMatchActionRequest(req dto.GetProfilesByMatchActionRequest) *userpbv1.GetProfilesByMatchActionRequest {
	return &userpbv1.GetProfilesByMatchActionRequest{
		Action: req.Action,
		Limit:  req.Limit,
		Offset: req.Offset,
	}
}

func MapGetProfilesByMatchActionResponse(resp *userpbv1.GetProfilesByMatchActionResponse) *dto.GetProfilesByMatchActionResponse {
	profiles := make([]dto.UserProfileRecommendation, len(resp.Profiles))
	for i, profile := range resp.Profiles {
		profiles[i] = dto.UserProfileRecommendation{
			ID:                profile.Id,
			FullName:          profile.FullName,
			ProfilePictureURL: profile.ProfilePictureUrl,
			Age:               int(profile.Age),
			HeightCm:          int(profile.HeightCm),
			MaritalStatus:     profile.MaritalStatus,
			Profession:        profile.Profession,
			HomeDistrict:      profile.HomeDistrict,
		}
	}

	pagination := dto.PaginationInfo{
		TotalCount: resp.Pagination.TotalCount,
		Limit:      int(resp.Pagination.Limit),
		Offset:     int(resp.Pagination.Offset),
		HasMore:    resp.Pagination.HasMore,
	}

	return &dto.GetProfilesByMatchActionResponse{
		Profiles:   profiles,
		Pagination: pagination,
	}
}