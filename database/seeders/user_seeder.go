package seeders

import (
	"goravel/app/models"

	"github.com/goravel/framework/facades"
)

type UserSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *UserSeeder) Signature() string {
	return "UserSeeder"
}

// Run executes the seeder logic.
func (s *UserSeeder) Run() error {
	password, _ := facades.Hash().Make("password1234")
	facades.Orm().Query().Create(&models.User{
		Name:     "User1",
		Email:    "user1@saas.com",
		Password: password,
		Role:     models.RoleUser,
	})

	return nil
}
