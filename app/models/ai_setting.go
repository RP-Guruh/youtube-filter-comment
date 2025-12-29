package models

import "github.com/goravel/framework/database/orm"

type AiSetting struct {
	orm.Model
	Name   string `json:"name"`
	Status bool   `json:"status"`
}
