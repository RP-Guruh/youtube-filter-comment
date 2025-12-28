package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type QuarantinedComment struct {
	orm.Model
	VideoID         uint       `json:"video_id"`
	CommentID       string     `json:"comment_id"`
	AuthorName      string     `json:"author_name"`
	CommentText     string     `json:"comment_text"`
	CommentCategory string     `json:"comment_category"`
	Status          string     `json:"status"`
	DeletedAt       *time.Time `json:"deleted_at" orm:"column:deleted_at"`
	DeletedBy       *uint      `json:"deleted_by"`
}
