package identity

import "context"

type contextKey string

func (k contextKey) String() string {
	return "context key: " + string(k)
}

var (
	userKey = contextKey("user")
)

// WithUser adds the user to the request context.
func WithUser(ctx context.Context, userId uint, roles []string) context.Context {
	id := &Identity{
		UserId: userId,
		Roles:  roles,
	}
	return context.WithValue(ctx, userKey, id)
}

// FromContext retrieves the user from the context.
func FromContext(ctx context.Context) *Identity {
	if val, ok := ctx.Value(userKey).(*Identity); ok {
		return val
	}
	return &Identity{}
}
