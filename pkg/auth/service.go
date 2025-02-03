package auth

import (
	"context"

	"github.com/sks/kihocche/pkg/auth/authutil"
)

type AuthService interface {
	CurrentUserPrincipal(ctx context.Context) (authutil.Principal, error)
}

func NewAuthService() AuthService {
	return authService{}
}

type authService struct{}

func (a authService) CurrentUserPrincipal(ctx context.Context) (authutil.Principal, error) {
	return authutil.Principal{}, nil
}
