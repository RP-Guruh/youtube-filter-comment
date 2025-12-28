package commands

import (
	"context"
	"goravel/app/helpers"
	"goravel/app/models"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type CheckForPendingScans struct {
}

// Signature The name and signature of the console command.
func (r *CheckForPendingScans) Signature() string {
	return "scan:check"
}

// Description The console command description.
func (r *CheckForPendingScans) Description() string {
	return "scan otomatis hanya untuk video yang di setting auto"
}

// Extend The console command extend.
func (r *CheckForPendingScans) Extend() command.Extend {
	return command.Extend{Category: "app"}
}

// Handle Execute the console command.
func (r *CheckForPendingScans) Handle(ctx console.Context) error {
	now := time.Now()

	// 1. Ambil semua VideoSettings dengan mode auto dan next_scan <= now
	var settings []models.VideoSetting
	err := facades.Orm().Query().Where("scan_mode", "auto").Where("next_scan <= ?", now).Get(&settings)
	if err != nil {
		return err
	}

	if len(settings) == 0 {
		return nil
	}
	log.Println("setting di command :", settings)
	// Ambil semua badwords yang aktif
	var badWords []models.BadWord
	_ = facades.Orm().Query().Where("is_active", true).Get(&badWords)
	normalizer := helpers.NewTextNormalize()

	for _, setting := range settings {
		// 2. Ambil info video
		var video models.Video
		_ = facades.Orm().Query().Where("id", setting.VideoID).First(&video)
		if video.ID == 0 {
			continue
		}

		// 3. Ambil YoutubeChannel tokens
		var youtubeChannel models.YoutubeChannel
		err = facades.Orm().Query().Where("user_id", video.UserID).Where("is_active", true).First(&youtubeChannel)
		if err != nil || youtubeChannel.ID == 0 {
			continue
		}

		// 4. Setup YouTube Service
		token := &oauth2.Token{
			AccessToken:  youtubeChannel.AccessToken,
			RefreshToken: youtubeChannel.RefreshToken,
			Expiry:       youtubeChannel.ExpiresAt,
		}
		config := helpers.GetGoogleConfig()
		httpClient := config.Client(context.Background(), token)
		youtubeService, err := youtube.NewService(context.Background(), option.WithHTTPClient(httpClient))
		if err != nil {
			log.Printf("[SCAN] Gagal membuat service untuk video %s: %v\n", video.VideoID, err)
			continue
		}

		// 5. Ambil 10 komentar terbaru
		call := youtubeService.CommentThreads.List([]string{"snippet"}).VideoId(video.VideoID).MaxResults(10).Order("time")
		response, err := call.Do()
		if err != nil {
			log.Printf("[SCAN] Gagal ambil komentar untuk video %s: %v\n", video.VideoID, err)
			continue
		}

		// 6. Proses komentar dengan Goroutines
		var wg sync.WaitGroup
		for _, item := range response.Items {
			wg.Add(1)
			go func(rawComment *youtube.CommentThread, v models.Video, s models.VideoSetting) {
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
					if s.ActionMode == "list" {
						count, _ := facades.Orm().Query().Model(&models.QuarantinedComment{}).Where("comment_id", commentID).Count()
						if count == 0 {
							qComment := models.QuarantinedComment{
								VideoID:         v.ID,
								CommentID:       commentID,
								AuthorName:      authorName,
								CommentText:     originalText,
								CommentCategory: category,
								Status:          "pending",
							}
							facades.Orm().Query().Create(&qComment)
						}
					} else if s.ActionMode == "deleted" {
						count, _ := facades.Orm().Query().Model(&models.LogComment{}).Where("comment_id", commentID).Count()
						if count == 0 {
							err := youtubeService.Comments.Delete(commentID).Do()
							if err != nil {
								log.Printf("[SCAN] Gagal hapus komentar %s: %v", commentID, err)
							} else {
								lComment := models.LogComment{
									VideoID:         v.ID,
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
			}(item, video, setting)
		}
		wg.Wait()

		// 7. Update jadwal scan
		lastScanned := time.Now()
		nextScan := lastScanned.Add(time.Minute * time.Duration(setting.FrequencyMinutes))

		setting.LastScanned = &lastScanned
		setting.NextScan = &nextScan
		facades.Orm().Query().Save(&setting)

		log.Printf("[SCAN] Video ID %d discan. Next: %v\n", video.ID, nextScan)
	}

	return nil
}
