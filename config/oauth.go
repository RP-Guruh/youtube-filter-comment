package config

import "github.com/goravel/framework/facades"

func init() {
	config := facades.Config()
	config.Add("oauth", map[string]any{
		"client_id":     config.Env("GOOGLE_CLIENT_ID", ""),
		"client_secret": config.Env("GOOGLE_CLIENT_SECRET", ""),
		"redirect_url":  config.Env("GOOGLE_REDIRECT_URI", "http://localhost:3000/api/auth/google/callback"),
	})
}
