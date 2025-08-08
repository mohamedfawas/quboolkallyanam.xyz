package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/client"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase/auth"
)

// Ensure fake implements the expanded client.AuthClient interface
var _ client.AuthClient = (*fakeAuthClient)(nil)

// fakeAuthClient implements only the behavior we need for UserRegister,
// while stubbing the rest to satisfy the interface.
type fakeAuthClient struct {
	resp *dto.UserRegisterResponse
	err  error

	// record last request for optional inspection
	lastReq dto.UserRegisterRequest
}

func (f *fakeAuthClient) UserRegister(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	f.lastReq = req
	return f.resp, f.err
}

func (f *fakeAuthClient) UserVerification(ctx context.Context, req dto.UserVerificationRequest) (*dto.UserVerificationResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeAuthClient) UserLogin(ctx context.Context, req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeAuthClient) UserLogout(ctx context.Context, accessToken string) error {
	return errors.New("not implemented")
}

func (f *fakeAuthClient) UserDelete(ctx context.Context, req dto.UserDeleteRequest) error {
	return errors.New("not implemented")
}

func (f *fakeAuthClient) AdminLogin(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminLoginResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeAuthClient) AdminLogout(ctx context.Context, accessToken string) error {
	return errors.New("not implemented")
}

func (f *fakeAuthClient) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeAuthClient) BlockOrUnblockUser(ctx context.Context, req dto.BlockOrUnblockUserRequest) (*dto.BlockOrUnblockUserResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeAuthClient) GetUsers(ctx context.Context, req dto.GetUsersRequest) (*dto.GetUsersResponse, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeAuthClient) GetUserByField(ctx context.Context, req dto.GetUserByFieldRequest) (*dto.GetUserByFieldResponse, error) {
	return nil, errors.New("not implemented")
}

func TestUserRegister(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		req            dto.UserRegisterRequest
		clientResp     *dto.UserRegisterResponse
		clientErr      error
		wantErr        error
		wantClientCall bool
	}{
		{
			name:    "invalid email",
			req:     dto.UserRegisterRequest{Email: "bad@", Phone: "+919876543210", Password: "P@ssw0rd"},
			wantErr: apperrors.ErrInvalidEmail,
		},
		{
			name:    "invalid phone",
			req:     dto.UserRegisterRequest{Email: "foo@bar.com", Phone: "notaphone", Password: "P@ssw0rd"},
			wantErr: apperrors.ErrInvalidPhoneNumber,
		},
		{
			name:    "invalid password",
			req:     dto.UserRegisterRequest{Email: "foo@bar.com", Phone: "+919876543210", Password: "weak"},
			wantErr: apperrors.ErrInvalidPassword,
		},
		{
			name:           "client error",
			req:            dto.UserRegisterRequest{Email: "fawas@gmail.com", Phone: "+919876543210", Password: "S3cur3Pwd!"},
			clientErr:      errors.New("rpc failure"),
			wantErr:        errors.New("rpc failure"),
			wantClientCall: true,
		},
		{
			name:           "success",
			req:            dto.UserRegisterRequest{Email: "fawas@gmail.com", Phone: "+919876543210", Password: "S3cur3Pwd!"},
			clientResp:     &dto.UserRegisterResponse{Email: "fawas@gmail.com", Phone: "+919876543210"},
			wantErr:        nil,
			wantClientCall: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := &fakeAuthClient{
				resp: tc.clientResp,
				err:  tc.clientErr,
			}
			uc := auth.NewAuthUsecase(fakeClient)

			got, err := uc.UserRegister(ctx, tc.req)

			// error assertions
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error(), "expected error")
				assert.Nil(t, got, "when error, response should be nil")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.clientResp, got, "response should match clientResp")
			}

			// ensure we only call the client on valid input
			if tc.wantClientCall {
				assert.Equal(t, tc.req, fakeClient.lastReq, "authClient.UserRegister must be called with same request")
			} else {
				assert.Zero(t, fakeClient.lastReq, "authClient should not be called on invalid input")
			}
		})
	}
}