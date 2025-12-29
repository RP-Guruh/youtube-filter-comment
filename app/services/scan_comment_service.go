package services

import (
	"context"
	"goravel/app/helpers"
	"goravel/app/models"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/goravel/framework/facades"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type ScanCommentService struct {
}

func NewScanCommentService() *ScanCommentService {
	return &ScanCommentService{}
}

type ScanResult struct {
	CommentID   string
	AuthorName  string
	CommentText string
	PublishedAt string
	IsBad       bool
	Category    string
}

func (s *ScanCommentService) ScanAndProcess(ctx context.Context, video models.Video, setting models.VideoSetting) ([]ScanResult, error) {
	// 1. Ambil YoutubeChannel tokens
	var youtubeChannel models.YoutubeChannel
	err := facades.Orm().Query().Where("user_id", video.UserID).Where("is_active", true).First(&youtubeChannel)
	if err != nil || youtubeChannel.ID == 0 {
		return nil, err
	}

	// 2. Setup YouTube Service
	token := &oauth2.Token{
		AccessToken:  youtubeChannel.AccessToken,
		RefreshToken: youtubeChannel.RefreshToken,
		Expiry:       youtubeChannel.ExpiresAt,
	}
	config := helpers.GetGoogleConfig()
	httpClient := config.Client(ctx, token)
	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	// 3. Ambil 10 komentar terbaru
	call := youtubeService.CommentThreads.List([]string{"snippet"}).VideoId(video.VideoID).MaxResults(10).Order("time")
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	// 4. Ambil badwords
	var badWords []models.BadWord
	err = facades.Orm().Query().Where("is_active", true).Get(&badWords)
	if err != nil {
		log.Println("[SERVICE] Gagal mengambil bad words:", err)
	}

	normalizer := helpers.NewTextNormalize()
	commentsToProcess := make(map[string]string)
	for _, item := range response.Items {
		snippet := item.Snippet.TopLevelComment.Snippet
		commentID := item.Snippet.TopLevelComment.Id
		commentsToProcess[commentID] = snippet.TextDisplay
	}

	// 5. Cek dengan Gemini AI (Hanya jika AiSetting ID 1 aktif)
	var aiResults map[string]bool
	var aiSetting models.AiSetting
	facades.Orm().Query().Where("id", 1).First(&aiSetting)

	if aiSetting.Status {
		aiResults = helpers.CheckJudolComments(commentsToProcess)
		log.Printf("[SERVICE] AI Check performed. Status: Active")
	} else {
		log.Printf("[SERVICE] AI Check skipped. Status: Inactive")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []ScanResult

	for _, item := range response.Items {
		wg.Add(1)
		go func(item *youtube.CommentThread) {
			defer wg.Done()

			snippet := item.Snippet.TopLevelComment.Snippet
			commentID := item.Snippet.TopLevelComment.Id
			originalText := snippet.TextDisplay
			authorName := snippet.AuthorDisplayName
			normalizedText := normalizer.Normalize(originalText)

			isBad := false
			category := ""

			// Check by bad words
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

			// Check Gemini AI if not detected yet
			if !isBad {
				if isJudol, ok := aiResults[commentID]; ok && isJudol {
					isBad = true
					category = "AI: Judi Online"
				}
			}

			if isBad {
				s.TakeAction(video, setting, youtubeService, commentID, authorName, originalText, category)
			}

			mu.Lock()
			results = append(results, ScanResult{
				CommentID:   commentID,
				AuthorName:  authorName,
				CommentText: originalText,
				PublishedAt: snippet.PublishedAt,
				IsBad:       isBad,
				Category:    category,
			})
			mu.Unlock()
		}(item)
	}
	wg.Wait()

	return results, nil
}

func (s *ScanCommentService) TakeAction(video models.Video, setting models.VideoSetting, youtubeService *youtube.Service, commentID, authorName, originalText, category string) {
	if setting.ActionMode == "list" {
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
			facades.Orm().Query().Create(&qComment)
		}
	} else if setting.ActionMode == "deleted" {
		count, _ := facades.Orm().Query().Model(&models.LogComment{}).Where("comment_id", commentID).Count()
		if count == 0 {
			err := youtubeService.Comments.Delete(commentID).Do()
			if err != nil {
				log.Printf("[SERVICE] Gagal hapus komentar %s: %v", commentID, err)
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
