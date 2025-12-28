package controllers

import (
	"goravel/app/helpers"
	"goravel/app/http/requests"
	"goravel/app/models"
	"log"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type SettingVideosController struct {
	// Dependent services
}

func NewSettingVideosController() *SettingVideosController {
	return &SettingVideosController{
		// Inject services
	}
}

func (r *SettingVideosController) Index(ctx http.Context) http.Response {
	user := helpers.UserInfo(ctx)
	if user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "User tidak terautentikasi",
		})
	}

	var settingVideos []models.VideoSetting
	err := facades.Orm().Query().
		Join("JOIN videos ON videos.id = video_settings.video_id").
		Where("videos.user_id", user.ID).
		Get(&settingVideos)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil data setting video: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, settingVideos)
}

func (r *SettingVideosController) Show(ctx http.Context) http.Response {
	user := helpers.UserInfo(ctx)
	if user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "User tidak terautentikasi",
		})
	}

	id := ctx.Request().Route("id")
	var settingVideo models.VideoSetting
	err := facades.Orm().Query().
		Join("JOIN videos ON videos.id = video_settings.video_id").
		Where("video_settings.id", id).
		Where("videos.user_id", user.ID).
		First(&settingVideo)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil data setting video: " + err.Error(),
		})
	}

	if settingVideo.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Setting video tidak ditemukan atau Anda tidak memiliki akses",
		})
	}

	return ctx.Response().Json(http.StatusOK, settingVideo)
}

func (r *SettingVideosController) Store(ctx http.Context) http.Response {
	// 1. validasi user
	user := helpers.UserInfo(ctx)
	if user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "User tidak terautentikasi",
		})
	}

	// 2. validasi request (post)
	var settingVideoRequest requests.SettingVideoRequest
	errors, err := ctx.Request().ValidateRequest(&settingVideoRequest)
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

	// 3. validasi video (pastikan milik user)
	var video models.Video
	facades.Orm().Query().Where("id", settingVideoRequest.VideoID).Where("user_id", user.ID).First(&video)
	if video.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Video tidak ditemukan atau Anda tidak memiliki akses",
		})
	}

	// 4. SIMPAN KE DB
	// next scan di ambil dari perhitungan antara frequency dan last scan

	var settingVideo models.VideoSetting
	settingVideo.VideoID = settingVideoRequest.VideoID
	settingVideo.ScanMode = settingVideoRequest.ScanMode
	settingVideo.FrequencyMinutes = settingVideoRequest.FrequencyMinutes
	settingVideo.ActionMode = settingVideoRequest.ActionMode

	if settingVideoRequest.ScanMode == "auto" {
		last_scan_schedule := time.Now()
		next_scan_schedule := last_scan_schedule.Add(time.Minute * time.Duration(settingVideoRequest.FrequencyMinutes))
		settingVideo.LastScanned = &last_scan_schedule
		settingVideo.NextScan = &next_scan_schedule
	} else {
		settingVideo.LastScanned = nil
		settingVideo.NextScan = nil
	}

	if err := facades.Orm().Query().UpdateOrCreate(&settingVideo, "video_id", settingVideoRequest.VideoID); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menyimpan setting video ke database: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message":       "Setting video berhasil disimpan",
		"setting_video": settingVideo,
	})
}

func (r *SettingVideosController) Update(ctx http.Context) http.Response {
	user := helpers.UserInfo(ctx)
	if user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "User tidak terautentikasi",
		})
	}

	id := ctx.Request().Route("id")
	var settingVideo models.VideoSetting
	err := facades.Orm().Query().
		Join("JOIN videos ON videos.id = video_settings.video_id").
		Where("video_settings.id", id).
		Where("videos.user_id", user.ID).
		First(&settingVideo)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mencek data setting video: " + err.Error(),
		})
	}

	if settingVideo.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Setting video tidak ditemukan atau Anda tidak memiliki akses",
		})
	}

	// 1. validasi request
	var settingVideoRequest requests.SettingVideoRequest
	errors, err := ctx.Request().ValidateRequest(&settingVideoRequest)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Gagal memproses request",
			"error":   err.Error(),
		})
	}
	if errors != nil && len(errors.All()) > 0 {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"errors": errors.All(),
		})
	}

	// 2. Update data
	settingVideo.ScanMode = settingVideoRequest.ScanMode
	settingVideo.FrequencyMinutes = settingVideoRequest.FrequencyMinutes
	settingVideo.ActionMode = settingVideoRequest.ActionMode

	if settingVideoRequest.ScanMode == "auto" {
		last_scan_schedule := time.Now()
		next_scan_schedule := last_scan_schedule.Add(time.Minute * time.Duration(settingVideoRequest.FrequencyMinutes))

		settingVideo.LastScanned = &last_scan_schedule
		settingVideo.NextScan = &next_scan_schedule
	} else {
		settingVideo.LastScanned = nil
		settingVideo.NextScan = nil
	}

	if err := facades.Orm().Query().Save(&settingVideo); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengupdate setting video: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message":       "Setting video berhasil diupdate",
		"setting_video": settingVideo,
	})
}

func (r *SettingVideosController) Destroy(ctx http.Context) http.Response {
	user := helpers.UserInfo(ctx)
	if user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "User tidak terautentikasi",
		})
	}

	id := ctx.Request().Route("id")
	log.Println("ID SETTING : ", id)
	var settingVideo models.VideoSetting
	err := facades.Orm().Query().
		Join("JOIN videos ON videos.id = video_settings.video_id").
		Where("video_settings.id", id).
		Where("videos.user_id", user.ID).
		First(&settingVideo)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mencek data setting video: " + err.Error(),
		})
	}

	if settingVideo.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Setting video tidak ditemukan atau Anda tidak memiliki akses",
		})
	}

	if _, err := facades.Orm().Query().Delete(&settingVideo); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menghapus setting video: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Setting video berhasil dihapus",
	})
}
