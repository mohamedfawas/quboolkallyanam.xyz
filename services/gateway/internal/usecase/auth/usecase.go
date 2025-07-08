package auth

import "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"

type AuthUsecase struct {
	authClient client.AuthClient
}

func NewAuthUsecase(authClient client.AuthClient) *AuthUsecase {
	return &AuthUsecase{authClient: authClient}
}
