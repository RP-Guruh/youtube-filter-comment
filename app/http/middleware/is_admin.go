package middleware

import (
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type IsAdmin struct {
}

func (r *IsAdmin) Handle(ctx http.Context, next http.Context) http.Response {
	var user models.User

	err := facades.Auth(ctx).User(&user)

	if err != nil || user.Role != "admin" {
		return ctx.Response().Json(http.StatusForbidden, map[string]string{
			"message": "Akses ditolak. Anda bukan admin.",
		})
	}

	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Akses diizinkan. Anda adalah admin.",
	})
}
