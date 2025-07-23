package user

import "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase"

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}
