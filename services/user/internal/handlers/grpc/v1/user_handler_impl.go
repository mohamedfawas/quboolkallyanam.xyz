package v1

import (
	"go.uber.org/zap"

	userpbv1 "github.com/mohamedfawas/quboolkallyanam.xyz/api/proto/user/v1"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/usecase"
)

type UserHandler struct {
	userpbv1.UnimplementedUserServiceServer
	userProfileUsecase usecase.UserProfileUsecase
	matchMakingUsecase usecase.MatchMakingUsecase
	logger             *zap.Logger
}

func NewUserHandler(
	userProfileUsecase usecase.UserProfileUsecase,
	matchMakingUsecase usecase.MatchMakingUsecase,
	logger *zap.Logger,
) *UserHandler {

	return &UserHandler{
		userProfileUsecase: userProfileUsecase,
		matchMakingUsecase: matchMakingUsecase,
		logger:             logger,
	}
}
