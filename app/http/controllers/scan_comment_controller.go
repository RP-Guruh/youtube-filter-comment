package controllers

import (
	"goravel/app/helpers"
	"goravel/app/models"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
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

	// 4. periksa melalui youtube data api v3 10 komentar terbaru

	// Ambil token dari youtube_channels
	var youtubeChannel models.YoutubeChannel
	err = facades.Orm().Query().Where("user_id", user.ID).Where("is_active", true).First(&youtubeChannel)
	if err != nil || youtubeChannel.ID == 0 {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"message": "Channel YouTube belum terhubung atau tidak aktif",
		})
	}

	// Buat token oauth2
	token := &oauth2.Token{
		AccessToken:  youtubeChannel.AccessToken,
		RefreshToken: youtubeChannel.RefreshToken,
		Expiry:       youtubeChannel.ExpiresAt,
	}

	config := helpers.GetGoogleConfig()
	client := config.Client(ctx.Context(), token)

	youtubeService, err := youtube.NewService(ctx.Context(), option.WithHTTPClient(client))
	if err != nil {
		log.Println("[ERROR] Gagal membuat YouTube service:", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal terhubung ke YouTube service",
		})
	}

	// Ambil 10 komentar terbaru
	call := youtubeService.CommentThreads.List([]string{"snippet"}).VideoId(video.VideoID).MaxResults(10).Order("time")
	response, err := call.Do()
	if err != nil {
		log.Println("[ERROR] Gagal mengambil komentar youtube:", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil komentar dari YouTube: " + err.Error(),
		})
	}

	// 5 & 6. pengecekan badwords (Goroutines) dan Aksi (Quarantine/Delete)

	// Ambil setting video
	var setting models.VideoSetting
	err = facades.Orm().Query().Where("video_id", video.ID).First(&setting)
	if err != nil || setting.ID == 0 {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Setting video tidak ditemukan",
		})
	}

	// Ambil semua badwords yang aktif
	var badWords []models.BadWord
	err = facades.Orm().Query().Where("is_active", true).Get(&badWords)
	if err != nil {
		log.Println("[ERROR] Gagal mengambil bad words:", err)
	}

	normalizer := helpers.NewTextNormalize()
	var wg sync.WaitGroup
	var mu sync.Mutex
	var filteredComments []http.Json

	for _, item := range response.Items {
		wg.Add(1)
		go func(rawComment *youtube.CommentThread) {
			defer wg.Done()

			snippet := rawComment.Snippet.TopLevelComment.Snippet
			commentID := rawComment.Snippet.TopLevelComment.Id
			originalText := snippet.TextDisplay
			authorName := snippet.AuthorDisplayName
			normalizedText := normalizer.Normalize(originalText)

			isBad := false
			category := ""

			for _, bw := range badWords {
				if bw.IsRegex {
					matched, _ := regexp.MatchString(bw.Word, normalizedText)
					if matched {
						isBad = true
						category = bw.Category
						break
					}
				} else {
					if strings.Contains(normalizedText, strings.ToLower(bw.Word)) {
						isBad = true
						category = bw.Category
						break
					}
				}
			}

			if isBad {
				if setting.ActionMode == "list" {
					log.Println("action mode :", setting.ActionMode)
					count, _ := facades.Orm().Query().Model(&models.QuarantinedComment{}).Where("comment_id", commentID).Count()
					if count == 0 {
						qComment := models.QuarantinedComment{
							VideoID:         video.ID,
							CommentID:       commentID,
							AuthorName:      authorName,
							CommentText:     originalText,
							CommentCategory: category,
							Status:          "pending",
						}
						err = facades.Orm().Query().Create(&qComment)
						if err != nil {
							log.Println("[ERROR] Gagal menambahkan komentar ke karantina:", err)
						}
					}
				} else if setting.ActionMode == "deleted" {
					count, _ := facades.Orm().Query().Model(&models.LogComment{}).Where("comment_text", originalText).Where("video_id", video.ID).Where("author_name", authorName).Count()

					if count == 0 {
						err := youtubeService.Comments.Delete(commentID).Do()
						if err != nil {
							log.Printf("[ERROR] Gagal hapus komentar %s: %v", commentID, err)
						} else {
							lComment := models.LogComment{
								VideoID:         video.ID,
								AuthorName:      authorName,
								CommentText:     originalText,
								CommentId:       commentID,
								CommentCategory: category,
								OriginAction:    "auto",
								FinalAction:     "deleted",
							}
							facades.Orm().Query().Create(&lComment)
						}
					}
				}
			}

			mu.Lock()
			filteredComments = append(filteredComments, http.Json{
				"comment_id":   commentID,
				"author_name":  authorName,
				"comment_text": originalText,
				"published_at": snippet.PublishedAt,
				"is_bad":       isBad,
				"category":     category,
			})
			mu.Unlock()
		}(item)
	}

	wg.Wait()

	return ctx.Response().Json(http.StatusOK, http.Json{
		"comments": filteredComments,
	})
}
