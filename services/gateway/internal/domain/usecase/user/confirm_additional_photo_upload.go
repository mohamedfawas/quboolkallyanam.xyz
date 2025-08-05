package user

import (
	"context"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) ConfirmAdditionalPhotoUpload(
	ctx context.Context, 
	req dto.ConfirmAdditionalPhotoUploadRequest) (*dto.ConfirmAdditionalPhotoUploadResponse, error) {
	
	resp, err := u.userClient.ConfirmAdditionalPhotoUpload(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}