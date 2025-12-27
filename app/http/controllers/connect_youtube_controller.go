package controllers

import (
	"context"
	"fmt"
	"goravel/app/helpers"
	"goravel/app/models"
	"log"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func getGoogleConfig() *oauth2.Config {
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

type ConnectYoutubeController struct {
	// Dependent services
}

func NewConnectYoutubeController() *ConnectYoutubeController {
	return &ConnectYoutubeController{
		// Inject services
	}
}

func (r *ConnectYoutubeController) Index(ctx http.Context) http.Response {

	user := helpers.UserInfo(ctx)
	config := getGoogleConfig()
	state := fmt.Sprintf("%d", user.ID)
	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return ctx.Response().Json(http.StatusOK, http.Json{
		"url": url,
	})

}

func (r *ConnectYoutubeController) YoutubeCallback(ctx http.Context) http.Response {

	code := ctx.Request().Query("code", "")
	userID := ctx.Request().Query("state", "")

	if code == "" {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"message": "Code tidak ditemukan"})
	}

	config := getGoogleConfig()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Println("[ERROR] Token Exchange Failed:", err.Error())
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"message": "Gagal tukar token",
			"error":   err.Error(),
		})
	}

	client := config.Client(context.Background(), token)
	youtubeService, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"message": "Gagal membuat YouTube service"})
	}

	call := youtubeService.Channels.List([]string{"snippet", "contentDetails"}).Mine(true)
	response, err := call.Do()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"message": "Gagal mengambil info channel YouTube: " + err.Error()})
	}

	if len(response.Items) == 0 {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"message": "Akun Google ini tidak memiliki YouTube Channel. Silakan buat channel terlebih dahulu.",
		})
	}

	channelItem := response.Items[0]
	channelID := channelItem.Id
	channelName := channelItem.Snippet.Title
	thumbnail := channelItem.Snippet.Thumbnails.Default.Url

	var youtubeChannel models.YoutubeChannel
	_ = facades.Orm().Query().Where("channel_id_youtube", channelID).First(&youtubeChannel)

	youtubeChannel.UserID = cast.ToUint(userID)
	youtubeChannel.ChannelIDYoutube = channelID
	youtubeChannel.ChannelName = channelName
	youtubeChannel.ChannelThumbnail = thumbnail
	youtubeChannel.AccessToken = token.AccessToken
	if token.RefreshToken != "" {
		youtubeChannel.RefreshToken = token.RefreshToken
	}
	youtubeChannel.ExpiresAt = token.Expiry
	youtubeChannel.IsActive = true

	err = facades.Orm().Query().Save(&youtubeChannel)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"message": "Gagal menyimpan data channel: " + err.Error()})
	}

	return ctx.Response().Json(200, map[string]any{
		"message": "Channel Berhasil Terhubung!",
		"channel": channelName,
	})

}
