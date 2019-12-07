package security

import "context"

type TokenProvider interface {
	SecurityContextFromToken(ctx context.Context, token string) (*UserContext, error)
	TokenFromSecurityContext(userContext *UserContext) (token string, err error)
}

var tokenProvider TokenProvider

func SetTokenProvider(provider TokenProvider) {
	if provider != nil {
		tokenProvider = provider
	}
}

func NewUserContextFromToken(ctx context.Context, token string) (userContext *UserContext, err error) {
	return tokenProvider.SecurityContextFromToken(ctx, token)
}