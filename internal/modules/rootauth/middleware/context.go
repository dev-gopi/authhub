package middleware

import (
	"context"

	"github.com/dev-gopi/authhub/internal/modules/rootauth/entity"
)

type contextKey string

const rootAuthContextKey contextKey = "root_auth_context"

func WithRootAuthContext(
	ctx context.Context,
	auth *entity.RootAuthContext,
) context.Context {
	return context.WithValue(
		ctx,
		rootAuthContextKey,
		auth,
	)
}

func FromRootAuthContext(
	ctx context.Context,
) (*entity.RootAuthContext, bool) {
	auth, ok := ctx.Value(
		rootAuthContextKey,
	).(*entity.RootAuthContext)

	return auth, ok
}
