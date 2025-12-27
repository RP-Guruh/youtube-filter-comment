package config

import "github.com/goravel/framework/facades"

func init() {
	config := facades.Config()
	config.Add("youtube", map[string]any{
		"youtube_key": config.Env("YOUTUBE_API_V3_KEY", ""),
	})
}
