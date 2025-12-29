package seeders

import (
	"goravel/app/models"
	"log"

	"github.com/goravel/framework/facades"
)

type AiSettingSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *AiSettingSeeder) Signature() string {
	return "AiSettingSeeder"
}

// Run executes the seeder logic.
func (s *AiSettingSeeder) Run() error {
	err := facades.Orm().Query().Create(&models.AiSetting{
		Name:   "gemini-2.5-flash-lite",
		Status: true,
	})
	log.Println(err)
	return nil
}
