package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type VideoRequest struct {
	UserID      uint   `json:"user_id" form:"user_id" orm:"user_id"`
	ChannelID   string `json:"channel_id" form:"channel_id" orm:"channel_id"`
	VideoID     string `json:"video_id" form:"video_id" orm:"video_id"`
	Title       string `json:"title" form:"title" orm:"title"`
	Description string `json:"description" form:"description" orm:"description"`
	Thumbnail   string `json:"thumbnail" form:"thumbnail" orm:"thumbnail"`
	PublishedAt string `json:"published_at" form:"published_at" orm:"published_at"`
}

func (r *VideoRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *VideoRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *VideoRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"user_id":      "required",
		"channel_id":   "required",
		"video_id":     "required",
		"title":        "required",
		"description":  "required",
		"thumbnail":    "required",
		"published_at": "required",
	}
}

func (r *VideoRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"user_id.required":      "User ID is required",
		"channel_id.required":   "Channel ID is required",
		"video_id.required":     "Video ID is required",
		"title.required":        "Title is required",
		"description.required":  "Description is required",
		"thumbnail.required":    "Thumbnail is required",
		"published_at.required": "Published at is required",
	}
}

func (r *VideoRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"user_id":      "User ID",
		"channel_id":   "Channel ID",
		"video_id":     "Video ID",
		"title":        "Title",
		"description":  "Description",
		"thumbnail":    "Thumbnail",
		"published_at": "Published at",
	}
}

func (r *VideoRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
