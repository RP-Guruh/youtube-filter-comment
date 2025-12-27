package controllers

import (
	"context"
	"goravel/app/helpers"
	"goravel/app/http/requests"
	"goravel/app/models"
	"log"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeController struct {
	// Dependent services
}

func NewYoutubeController() *YoutubeController {
	return &YoutubeController{
		// Inject services
	}
}

func (r *YoutubeController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *YoutubeController) ListVideo(ctx http.Context) http.Response {
	// cek youtube channel user
	user := helpers.UserInfo(ctx)
	if user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "User tidak terautentikasi",
		})
	}

	var youtubeChannel models.YoutubeChannel
	err := facades.Orm().Query().Where("user_id", user.ID).First(&youtubeChannel)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Youtube channel tidak ditemukan untuk user ini",
			"user_id": user.ID,
		})
	}

	config := helpers.GetGoogleConfig()
	token := &oauth2.Token{
		AccessToken:  youtubeChannel.AccessToken,
		RefreshToken: youtubeChannel.RefreshToken,
		Expiry:       youtubeChannel.ExpiresAt,
	}

	client := config.Client(context.Background(), token)
	youtubeService, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal membuat YouTube service",
		})
	}

	call := youtubeService.Search.List([]string{"snippet"}).
		ChannelId(youtubeChannel.ChannelIDYoutube).
		Type("video").
		MaxResults(50).
		Order("date")

	response, err := call.Do()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil list video: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"channel": youtubeChannel,
		"videos":  response.Items,
	})
}

func (r *YoutubeController) AddVideo(ctx http.Context) http.Response {
	// 1. Ambil user yang sedang login
	user := helpers.UserInfo(ctx)
	if user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "User tidak terautentikasi",
		})
	}

	// 2. Validate and Bind Request
	var videoRequest requests.VideoRequest
	errors, err := ctx.Request().ValidateRequest(&videoRequest)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Gagal memproses request (Format JSON salah atau link putus)",
			"error":   err.Error(),
		})
	}

	if errors != nil && len(errors.All()) > 0 {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"errors": errors.All(),
		})
	}

	log.Printf("[DEBUG] VideoRequest: UserID=%d, ChannelID=%s, VideoID=%s\n", videoRequest.UserID, videoRequest.ChannelID, videoRequest.VideoID)

	// 3. Pastikan YouTube Channel tersebut milik user yang sedang login
	var youtubeChannel models.YoutubeChannel
	err = facades.Orm().Query().Where("user_id", user.ID).Where("channel_id_youtube", videoRequest.ChannelID).First(&youtubeChannel)

	if err != nil || youtubeChannel.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Youtube channel tidak ditemukan atau bukan milik Anda",
		})
	}

	// 4. Proses Simpan Video
	var youtubeVideo models.Video
	facades.Orm().Query().Where("video_id", videoRequest.VideoID).First(&youtubeVideo)

	youtubeVideo.UserID = user.ID
	youtubeVideo.ChannelID = youtubeChannel.ChannelIDYoutube
	youtubeVideo.VideoID = videoRequest.VideoID
	youtubeVideo.Title = videoRequest.Title
	youtubeVideo.Description = videoRequest.Description
	youtubeVideo.Thumbnail = videoRequest.Thumbnail
	youtubeVideo.PublishedAt = videoRequest.PublishedAt

	if err := facades.Orm().Query().Save(&youtubeVideo); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menyimpan video ke database: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Video berhasil disimpan",
		"video":   youtubeVideo,
	})
}
