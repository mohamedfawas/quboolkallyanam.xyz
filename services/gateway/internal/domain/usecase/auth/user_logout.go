package auth

import (
	"context"
)

func (u *authUsecase) UserLogout(
	ctx context.Context, 
	accessToken string) error {
		
	return u.authClient.UserLogout(ctx, accessToken)
}
