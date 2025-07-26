package user

import (
	"time"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

////////////////////////////// Update User Profile //////////////////////////////

func MapUpdateUserProfileRequest(req dto.UserProfilePatchRequest) *userpbv1.UpdateUserProfileRequest {
	grpcReq := &userpbv1.UpdateUserProfileRequest{}

	grpcReq.IsBride = &wrapperspb.BoolValue{Value: req.IsBride}
	grpcReq.PhysicallyChallenged = &wrapperspb.BoolValue{Value: req.PhysicallyChallenged}

	if req.FullName != nil {
		grpcReq.FullName = &wrapperspb.StringValue{Value: *req.FullName}
	}

	if req.DateOfBirth != nil {
		// Convert date from "2006-01-02" format to RFC3339 format
		date, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err == nil {
			rfc3339Date := date.Format(time.RFC3339)
			grpcReq.DateOfBirth = &wrapperspb.StringValue{Value: rfc3339Date}
		} else {
			// If parsing fails, pass the original value (let the user service handle the error)
			grpcReq.DateOfBirth = &wrapperspb.StringValue{Value: *req.DateOfBirth}
		}
	}
	if req.Community != nil {
		grpcReq.Community = &wrapperspb.StringValue{Value: *req.Community}
	}
	if req.MaritalStatus != nil {
		grpcReq.MaritalStatus = &wrapperspb.StringValue{Value: *req.MaritalStatus}
	}
	if req.Profession != nil {
		grpcReq.Profession = &wrapperspb.StringValue{Value: *req.Profession}
	}
	if req.ProfessionType != nil {
		grpcReq.ProfessionType = &wrapperspb.StringValue{Value: *req.ProfessionType}
	}
	if req.HighestEducationLevel != nil {
		grpcReq.HighestEducationLevel = &wrapperspb.StringValue{Value: *req.HighestEducationLevel}
	}
	if req.HomeDistrict != nil {
		grpcReq.HomeDistrict = &wrapperspb.StringValue{Value: *req.HomeDistrict}
	}

	if req.HeightCm != nil {
		grpcReq.HeightCm = &wrapperspb.Int32Value{Value: int32(*req.HeightCm)}
	}

	return grpcReq
}

func MapUpdateUserProfileResponse(resp *userpbv1.UpdateUserProfileResponse) *dto.UpdateUserProfileResponse {
	return &dto.UpdateUserProfileResponse{
		Success: resp.Success.GetValue(),
	}
}

////////////////////////////// Update User Partner Preferences //////////////////////////////

func MapUpdateUserPartnerPreferencesRequest(operationType string, req dto.PartnerPreferencePatchRequest) *userpbv1.UpdateUserPartnerPreferencesRequest {
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
		grpcReq.PreferredCommunities = *req.PreferredCommunities
	}

	if req.PreferredMaritalStatus != nil {
		grpcReq.PreferredMaritalStatus = *req.PreferredMaritalStatus
	}

	if req.PreferredProfessions != nil {
		grpcReq.PreferredProfessions = *req.PreferredProfessions
	}

	if req.PreferredProfessionTypes != nil {
		grpcReq.PreferredProfessionTypes = *req.PreferredProfessionTypes
	}

	if req.PreferredEducationLevels != nil {
		grpcReq.PreferredEducationLevels = *req.PreferredEducationLevels
	}

	if req.PreferredHomeDistricts != nil {
		grpcReq.PreferredHomeDistricts = *req.PreferredHomeDistricts
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
		TargetProfileId: uint32(req.TargetProfileID),
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