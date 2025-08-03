package user

import (
	"context"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

func (u *userUsecase) ConfirmProfilePhotoUpload(ctx context.Context, req dto.ConfirmProfilePhotoUploadRequest) (*dto.ConfirmProfilePhotoUploadResponse, error) {
	resp, err := u.userClient.ConfirmProfilePhotoUpload(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}