package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// OpenApiAuth
// check http query params: access_token
// check http header Authorization: Bearer *** or header x-access-token
func OpenApiAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {

	}
}
