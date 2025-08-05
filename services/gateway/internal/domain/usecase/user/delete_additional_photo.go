package user

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) DeleteAdditionalPhoto(
	ctx context.Context, 
	req dto.DeleteAdditionalPhotoRequest) (*dto.DeleteAdditionalPhotoResponse, error) {
	
	if req.DisplayOrder < 1 || req.DisplayOrder > 3 {
		return nil, apperrors.ErrInvalidDisplayOrder
	}
	
	resp, err := u.userClient.DeleteAdditionalPhoto(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}