package requests

import (
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type SettingVideoRequest struct {
	VideoID          uint      `json:"video_id" orm:"video_id" form:"video_id"`
	ScanMode         string    `json:"scan_mode" orm:"scan_mode" form:"scan_mode"`
	FrequencyMinutes int16     `json:"frequency_minutes" orm:"frequency_minutes" form:"frequency_minutes"`
	ActionMode       string    `json:"action_mode" orm:"action_mode" form:"action_mode"`
	LastScanned      time.Time `json:"last_scanned" orm:"last_scanned" form:"last_scanned"`
	NextScan         time.Time `json:"next_scan" orm:"next_scan" form:"next_scan"`
}

func (r *SettingVideoRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *SettingVideoRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *SettingVideoRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"video_id":          "required",
		"scan_mode":         "required",
		"frequency_minutes": "required",
		"action_mode":       "required",
	}
}

func (r *SettingVideoRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"video_id.required":          "ID video wajib diisi",
		"scan_mode.required":         "Mode scan wajib diisi",
		"frequency_minutes.required": "Frekuensi menit wajib diisi",
		"action_mode.required":       "Mode tindakan wajib diisi",
	}
}

func (r *SettingVideoRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"video_id":          "ID video",
		"scan_mode":         "Mode scan",
		"frequency_minutes": "Frekuensi menit",
		"action_mode":       "Mode tindakan",
	}
}

func (r *SettingVideoRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
