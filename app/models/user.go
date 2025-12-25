package models

import (
	"github.com/goravel/framework/database/orm"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	orm.Model
	Name     string
	Email    string
	Password string `json:"-"`
	Role     string
}
