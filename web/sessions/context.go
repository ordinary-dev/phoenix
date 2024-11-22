package sessions

import (
	"context"
	"net/http"
)

type MyContextKey string

const (
	UsernameKey MyContextKey = "username"
)

func AddUsernameToContext(r *http.Request, username string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, UsernameKey, username)
	return r.WithContext(ctx)
}

func GetUsername(ctx context.Context) string {
	val := ctx.Value(UsernameKey)
	return val.(string)
}
