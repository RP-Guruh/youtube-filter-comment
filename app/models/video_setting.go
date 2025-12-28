package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type VideoSetting struct {
	orm.Model
	VideoID          uint       `json:"video_id" orm:"video_id"`
	ScanMode         string     `json:"scan_mode" orm:"scan_mode"`
	FrequencyMinutes int16      `json:"frequency_minutes" orm:"frequency_minutes"`
	ActionMode       string     `json:"action_mode" orm:"action_mode"`
	LastScanned      *time.Time `json:"last_scanned" orm:"last_scanned"`
	NextScan         *time.Time `json:"next_scan" orm:"next_scan"`
}
