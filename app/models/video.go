package models

import (
	"github.com/goravel/framework/database/orm"
)

type Video struct {
	orm.Model
	UserID      uint   `json:"user_id" orm:"user_id"`
	ChannelID   string `json:"channel_id" orm:"channel_id"`
	VideoID     string `json:"video_id" orm:"video_id"`
	Title       string `json:"title" orm:"title"`
	Description string `json:"description" orm:"description"`
	Thumbnail   string `json:"thumbnail" orm:"thumbnail"`
	PublishedAt string `json:"published_at" orm:"published_at"`
}
