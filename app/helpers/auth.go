package helpers

import (
	"goravel/app/models"
	"log"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func UserInfo(ctx http.Context) models.User {
	var user models.User

	if guard := ctx.Request().Header("Guard"); guard == "" {
		_ = facades.Auth(ctx).User(&user)
	} else {
		_ = facades.Auth(ctx).Guard(guard).User(&user)
	}

	return user
}

func GetGoogleConfig() *oauth2.Config {
	clientId := facades.Config().GetString("oauth.client_id")
	redirectUrl := facades.Config().GetString("oauth.redirect_url")

	log.Printf("[DEBUG] OAuth Config - ClientID: %s..., RedirectURL: %s", clientId[:10], redirectUrl)

	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: facades.Config().GetString("oauth.client_secret"),
		RedirectURL:  redirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/youtube.force-ssl",
		},
		Endpoint: google.Endpoint,
	}
}
