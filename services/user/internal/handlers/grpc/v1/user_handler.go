package v1

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/contextutils"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase"
)

type UserHandler struct {
	userpbv1.UnimplementedUserServiceServer
	userProfileUsecase usecase.UserProfileUsecase
	matchMakingUsecase usecase.MatchMakingUsecase
}

func NewUserHandler(
	userProfileUsecase usecase.UserProfileUsecase,
	matchMakingUsecase usecase.MatchMakingUsecase,
) *UserHandler {

	return &UserHandler{
		userProfileUsecase: userProfileUsecase,
		matchMakingUsecase: matchMakingUsecase,
	}
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *userpbv1.UpdateUserProfileRequest) (*userpbv1.UpdateUserProfileResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	entityReq := entity.UpdateUserProfileRequest{}
	if req.IsBride != nil {
		entityReq.IsBride = &req.IsBride.Value
	}
	if req.FullName != nil {
		entityReq.FullName = &req.FullName.Value
	}
	if req.DateOfBirth != nil {
		dateOfBirth, err := time.Parse(time.RFC3339, req.DateOfBirth.Value)
		if err != nil {
			log.Printf("Failed to parse date of birth: %v", err)
			return nil, status.Errorf(codes.InvalidArgument, "invalid date of birth: %v", err)
		}
		entityReq.DateOfBirth = &dateOfBirth
	}
	if req.HeightCm != nil {
		height := int(req.HeightCm.Value)
		entityReq.HeightCm = &height
	}
	if req.PhysicallyChallenged != nil {
		entityReq.PhysicallyChallenged = &req.PhysicallyChallenged.Value
	}
	if req.Community != nil {
		entityReq.Community = &req.Community.Value
	}
	if req.MaritalStatus != nil {
		entityReq.MaritalStatus = &req.MaritalStatus.Value
	}
	if req.Profession != nil {
		entityReq.Profession = &req.Profession.Value
	}
	if req.ProfessionType != nil {
		entityReq.ProfessionType = &req.ProfessionType.Value
	}
	if req.HighestEducationLevel != nil {
		entityReq.HighestEducationLevel = &req.HighestEducationLevel.Value
	}
	if req.HomeDistrict != nil {
		entityReq.HomeDistrict = &req.HomeDistrict.Value
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	err = h.userProfileUsecase.UpdateUserProfile(ctx, userIDUUID, entityReq)
	if err != nil {
		log.Printf("Failed to update user profile: %v", err)
		return nil, err
	}
	return &userpbv1.UpdateUserProfileResponse{Success: &wrapperspb.BoolValue{Value: true}}, nil
}

func (h *UserHandler) UpdateUserPartnerPreferences(ctx context.Context, req *userpbv1.UpdateUserPartnerPreferencesRequest) (*userpbv1.UpdateUserPartnerPreferencesResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	operationType := ""
	if req.OperationType != nil {
		operationType = req.OperationType.Value
	}

	entityReq := entity.UpdateUserPartnerPreferencesRequest{}

	if req.MinAgeYears != nil {
		minAge := int(req.MinAgeYears.Value)
		entityReq.MinAgeYears = &minAge
	}
	if req.MaxAgeYears != nil {
		maxAge := int(req.MaxAgeYears.Value)
		entityReq.MaxAgeYears = &maxAge
	}
	if req.MinHeightCm != nil {
		minHeight := int(req.MinHeightCm.Value)
		entityReq.MinHeightCM = &minHeight
	}
	if req.MaxHeightCm != nil {
		maxHeight := int(req.MaxHeightCm.Value)
		entityReq.MaxHeightCM = &maxHeight
	}
	if req.AcceptPhysicallyChallenged != nil {
		entityReq.AcceptPhysicallyChallenged = &req.AcceptPhysicallyChallenged.Value
	}

	if len(req.PreferredCommunities) > 0 {
		entityReq.PreferredCommunities = &req.PreferredCommunities
	}
	if len(req.PreferredMaritalStatus) > 0 {
		entityReq.PreferredMaritalStatus = &req.PreferredMaritalStatus
	}
	if len(req.PreferredProfessions) > 0 {
		entityReq.PreferredProfessions = &req.PreferredProfessions
	}
	if len(req.PreferredProfessionTypes) > 0 {
		entityReq.PreferredProfessionTypes = &req.PreferredProfessionTypes
	}
	if len(req.PreferredEducationLevels) > 0 {
		entityReq.PreferredEducationLevels = &req.PreferredEducationLevels
	}
	if len(req.PreferredHomeDistricts) > 0 {
		entityReq.PreferredHomeDistricts = &req.PreferredHomeDistricts
	}

	err = h.userProfileUsecase.UpdateUserPartnerPreferences(ctx, userIDUUID, operationType, entityReq)
	if err != nil {
		log.Printf("Failed to update partner preferences: %v", err)
		return nil, err
	}

	return &userpbv1.UpdateUserPartnerPreferencesResponse{Success: &wrapperspb.BoolValue{Value: true}}, nil
}

func (h *UserHandler) RecordMatchAction(ctx context.Context, req *userpbv1.RecordMatchActionRequest) (*userpbv1.RecordMatchActionResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	if req.Action == "" {
		log.Printf("Action is required")
		return nil, status.Errorf(codes.InvalidArgument, "action is required")
	}

	if req.TargetProfileId == 0 {
		log.Printf("Target profile ID is required")
		return nil, status.Errorf(codes.InvalidArgument, "target profile ID is required")
	}

	success, err := h.matchMakingUsecase.RecordMatchAction(ctx, userIDUUID, uint(req.TargetProfileId), req.Action)
	if err != nil {
		log.Printf("Failed to record match action: %v", err)
		return nil, err
	}

	return &userpbv1.RecordMatchActionResponse{Success: success}, nil
}

func (h *UserHandler) GetMatchRecommendations(ctx context.Context, req *userpbv1.GetMatchRecommendationsRequest) (*userpbv1.GetMatchRecommendationsResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	profiles, pagination, err := h.matchMakingUsecase.RecommendUserProfiles(ctx, userIDUUID, limit, offset)
	if err != nil {
		log.Printf("Failed to get match recommendations: %v", err)
		return nil, err
	}

	protoProfiles := make([]*userpbv1.UserProfileRecommendation, len(profiles))
	for i, profile := range profiles {
		protoProfiles[i] = &userpbv1.UserProfileRecommendation{
			Id:                profile.ID,
			FullName:          profile.FullName,
			ProfilePictureUrl: getStringValue(profile.ProfilePictureURL),
			Age:               int32(profile.Age),
			HeightCm:          int32(profile.HeightCm),
			MaritalStatus:     profile.MaritalStatus,
			Profession:        profile.Profession,
			HomeDistrict:      profile.HomeDistrict,
		}
	}

	protoPagination := &userpbv1.PaginationInfo{
		TotalCount: pagination.TotalCount,
		Limit:      int32(pagination.Limit),
		Offset:     int32(pagination.Offset),
		HasMore:    pagination.HasMore,
	}

	return &userpbv1.GetMatchRecommendationsResponse{
		Profiles:   protoProfiles,
		Pagination: protoPagination,
	}, nil
}

func (h *UserHandler) GetProfilesByMatchAction(ctx context.Context, req *userpbv1.GetProfilesByMatchActionRequest) (*userpbv1.GetProfilesByMatchActionResponse, error) {
	userID, err := contextutils.GetUserID(ctx)
	if err != nil {
		log.Printf("Failed to get user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "user ID not found: %v", err)
	}

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Failed to parse user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	if req.Action == "" {
		log.Printf("Action is required")
		return nil, status.Errorf(codes.InvalidArgument, "action is required")
	}

	limit := int(req.Limit)
	offset := int(req.Offset)

	profiles, pagination, err := h.matchMakingUsecase.GetProfilesByMatchAction(ctx, userIDUUID, req.Action, limit, offset)
	if err != nil {
		log.Printf("Failed to get profiles by match action: %v", err)
		return nil, err
	}

	protoProfiles := make([]*userpbv1.UserProfileRecommendation, len(profiles))
	for i, profile := range profiles {
		protoProfiles[i] = &userpbv1.UserProfileRecommendation{
			Id:                profile.ID,
			FullName:          profile.FullName,
			ProfilePictureUrl: getStringValue(profile.ProfilePictureURL),
			Age:               int32(profile.Age),
			HeightCm:          int32(profile.HeightCm),
			MaritalStatus:     profile.MaritalStatus,
			Profession:        profile.Profession,
			HomeDistrict:      profile.HomeDistrict,
		}
	}

	protoPagination := &userpbv1.PaginationInfo{
		TotalCount: pagination.TotalCount,
		Limit:      int32(pagination.Limit),
		Offset:     int32(pagination.Offset),
		HasMore:    pagination.HasMore,
	}

	return &userpbv1.GetProfilesByMatchActionResponse{
		Profiles:   protoProfiles,
		Pagination: protoPagination,
	}, nil
}

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
