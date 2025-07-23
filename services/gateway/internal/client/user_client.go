package client

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
)

type UserClient interface {
	UpdateUserProfile(ctx context.Context, req dto.UserProfilePatchRequest) error
	UpdateUserPartnerPreferences(ctx context.Context, operationType string, req dto.PartnerPreferencePatchRequest) error
	RecordMatchAction(ctx context.Context, req dto.RecordMatchActionRequest) (*dto.RecordMatchActionResponse, error)
}
