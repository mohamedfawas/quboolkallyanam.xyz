package auth

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
)

type authUsecase struct {
	authClient client.AuthClient
}

func NewAuthUsecase(authClient client.AuthClient) usecase.AuthUsecase {
	return &authUsecase{authClient: authClient}
}
