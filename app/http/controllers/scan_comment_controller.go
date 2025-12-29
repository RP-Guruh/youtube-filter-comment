package controllers

import (
	"goravel/app/helpers"
	"goravel/app/models"
	"goravel/app/services"
	"log"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ScanCommentController struct {
	// Dependent services
}

func NewScanCommentController() *ScanCommentController {
	return &ScanCommentController{
		// Inject services
	}
}

func (r *ScanCommentController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *ScanCommentController) ScanComment(ctx http.Context) http.Response {
	// 1. validasi user
	user := helpers.UserInfo(ctx)
	if user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "User tidak terautentikasi",
		})
	}

	// 2. tangkap param id
	id_video := ctx.Request().Route("id")

	// 3. periksa ke db berdasarkan id video, ambil video-id youtube nya
	var video models.Video
	err := facades.Orm().Query().Where("id", id_video).First(&video)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Video tidak ditemukan",
		})
	}
	if video.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Video tidak ditemukan",
		})
	}

	// 3. Ambil setting video
	var setting models.VideoSetting
	err = facades.Orm().Query().Where("video_id", video.ID).First(&setting)
	if err != nil || setting.ID == 0 {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Setting video tidak ditemukan",
		})
	}

	// 4. Hubungkan ke service baru
	scanService := services.NewScanCommentService()
	serviceResults, err := scanService.ScanAndProcess(ctx.Context(), video, setting)
	if err != nil {
		log.Println("[ERROR] Gagal proses scan melalui service:", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal melakukan scan: " + err.Error(),
		})
	}

	// Format results for response
	var filteredComments []http.Json
	for _, res := range serviceResults {
		filteredComments = append(filteredComments, http.Json{
			"comment_id":   res.CommentID,
			"author_name":  res.AuthorName,
			"comment_text": res.CommentText,
			"published_at": res.PublishedAt,
			"is_bad":       res.IsBad,
			"category":     res.Category,
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"comments": filteredComments,
	})
}
