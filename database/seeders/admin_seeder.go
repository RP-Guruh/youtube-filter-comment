package seeders

import (
	"goravel/app/models"

	"github.com/goravel/framework/facades"
)

type AdminSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *AdminSeeder) Signature() string {
	return "AdminSeeder"
}

// Run executes the seeder logic.
func (s *AdminSeeder) Run() error {
	password, _ := facades.Hash().Make("password1234")
	facades.Orm().Query().Create(&models.User{
		Name:     "Guruh Admin",
		Email:    "admin@saas.com",
		Password: password,
		Role:     models.RoleAdmin,
	})

	return nil
}
