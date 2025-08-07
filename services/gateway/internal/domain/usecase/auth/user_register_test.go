package auth_test

import (
    "context"
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/apperrors"
    "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/dto"
    "github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/domain/usecase/auth"
)

// fakeAuthClient implements only the method we need.
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
            clientResp:     &dto.UserRegisterResponse{UserID: "uuid-1234", Email: "fawas@gmail.com"},
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
            uc := auth.NewAuthUsecase(/* pass in fakeClient as the authClient */ fakeClient)

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
