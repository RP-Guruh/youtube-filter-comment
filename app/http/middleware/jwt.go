package middleware

import (
	"github.com/goravel/framework/contracts/http"
)

func Jwt() http.Middleware {
	return func(ctx http.Context) {
		ctx.Request().Next()
	}
}
