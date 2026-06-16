package middleware

import (
	"context"

	"Order-Management-System/internal/application/auth"
)

type actorContextKey struct{}

func WithActor(ctx context.Context, actor auth.Actor) context.Context {
	return context.WithValue(ctx, actorContextKey{}, actor)
}

func ActorFromContext(ctx context.Context) (auth.Actor, bool) {
	actor, ok := ctx.Value(actorContextKey{}).(auth.Actor)
	return actor, ok
}
