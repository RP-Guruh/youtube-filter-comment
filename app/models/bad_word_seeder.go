package models

import "github.com/goravel/framework/database/orm"

type BadWord struct {
	orm.Model
	Word          string  `json:"word"`
	Category      string  `json:"category"`
	IsRegex       bool    `json:"is_regex"`
	SeverityScore float64 `json:"severity_score"`
	IsActive      bool    `json:"is_active"`
}
