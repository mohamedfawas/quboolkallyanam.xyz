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

func (h *UserHandler) PutUserProfile(c *gin.Context) {
	var req dto.UserProfilePutRequest
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

	request := dto.UserProfilePatchRequest{
		IsBride:               req.IsBride,
		FullName:              req.FullName,
		DateOfBirth:           req.DateOfBirth,
		HeightCm:              req.HeightCm,
		PhysicallyChallenged:  req.PhysicallyChallenged,
		Community:             req.Community,
		MaritalStatus:         req.MaritalStatus,
		Profession:            req.Profession,
		ProfessionType:        req.ProfessionType,
		HighestEducationLevel: req.HighestEducationLevel,
		HomeDistrict:          req.HomeDistrict,
	}

	err := h.userUsecase.UpdateUserProfile(ctx, request)
	if err != nil {
		log.Printf("Failed to update user profile: %v", err)
		apiresponse.Fail(c, err)
		return
	}

	apiresponse.Success(c, "User profile updated successfully", nil)
}
