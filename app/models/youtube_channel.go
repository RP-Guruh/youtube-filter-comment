package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type YoutubeChannel struct {
	orm.Model
	UserID           uint
	ChannelIDYoutube string
	ChannelName      string
	ChannelThumbnail string
	AccessToken      string
	RefreshToken     string
	ExpiresAt        time.Time
	IsActive         bool
}
