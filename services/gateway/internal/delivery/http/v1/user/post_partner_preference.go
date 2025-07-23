package user

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
	apiresponse "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/apiresponse"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (h *UserHandler) PostPartnerPreference(c *gin.Context) {
	var req dto.PartnerPreferenceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		apiresponse.Fail(c, fmt.Errorf("invalid request body: %w", err))
		return
	}

	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		apiresponse.Fail(c, fmt.Errorf("user ID not found in context"))
		return
	}

	ctx := context.WithValue(c.Request.Context(), constants.ContextKeyUserID, userID)

	request := dto.PartnerPreferencePatchRequest{
		MinAgeYears:                req.MinAgeYears,
		MaxAgeYears:                req.MaxAgeYears,
		MinHeightCM:                req.MinHeightCM,
		MaxHeightCM:                req.MaxHeightCM,
		AcceptPhysicallyChallenged: &req.AcceptPhysicallyChallenged,
		PreferredCommunities:       &req.PreferredCommunities,
		PreferredMaritalStatus:     &req.PreferredMaritalStatus,
		PreferredProfessions:       &req.PreferredProfessions,
		PreferredProfessionTypes:   &req.PreferredProfessionTypes,
		PreferredEducationLevels:   &req.PreferredEducationLevels,
		PreferredHomeDistricts:     &req.PreferredHomeDistricts,
	}
	err := h.userUsecase.UpdateUserPartnerPreferences(ctx, constants.CreateOperationType, request)
	if err != nil {
		log.Printf("Failed to update partner preference: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "Partner preference updated successfully", nil)
}
