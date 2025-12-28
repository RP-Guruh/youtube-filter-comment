package models

import (
	"github.com/goravel/framework/database/orm"
)

type LogComment struct {
	orm.Model
	VideoID         uint   `json:"video_id"`
	AuthorName      string `json:"author_name"`
	CommentText     string `json:"comment_text"`
	CommentId       string `json:"comment_id"`
	CommentCategory string `json:"comment_category"`
	OriginAction    string `json:"origin_action"`
	FinalAction     string `json:"final_action"`
}
