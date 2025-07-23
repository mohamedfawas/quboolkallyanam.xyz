package user

import (
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"
)

type userUsecase struct {
	userClient client.UserClient
}

func NewUserUsecase(userClient client.UserClient) usecase.UserUsecase {
	return &userUsecase{userClient: userClient}
}
