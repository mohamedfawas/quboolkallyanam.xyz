package auth

import (
	"context"
)

func (u *authUsecase) AdminLogout(
	ctx context.Context,
	accessToken string) error {

	return u.authClient.AdminLogout(ctx, accessToken)
}
