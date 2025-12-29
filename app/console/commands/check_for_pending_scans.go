package commands

import (
	"context"
	"goravel/app/models"
	"goravel/app/services"
	"log"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
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

		// 4-6. Proses scan melalui service
		scanService := services.NewScanCommentService()
		_, err = scanService.ScanAndProcess(context.Background(), video, setting)
		if err != nil {
			log.Printf("[SCAN] Gagal proses video %s melalui service: %v\n", video.VideoID, err)
			continue
		}

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
