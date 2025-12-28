package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	"goravel/app/http/middleware"
)

func Api() {
	userController := controllers.NewUserController()
	authController := controllers.NewAuthController()
	connectyoutubeController := controllers.NewConnectYoutubeController()
	youtubeController := controllers.NewYoutubeController()
	settingVideosController := controllers.NewSettingVideosController()
	scanCommentController := controllers.NewScanCommentController()

	facades.Route().Get("/users/{id}", userController.Show)

	facades.Route().Prefix("api").Group(func(router route.Router) {
		// Auth
		router.Post("/auth/login", authController.Login)
		router.Middleware(middleware.Auth()).Get("/auth/info", authController.Info)

		// Youtube
		router.Middleware(middleware.Auth()).Get("/connect-youtube", connectyoutubeController.Index)
		router.Get("/auth/google/callback", connectyoutubeController.YoutubeCallback)
		router.Middleware(middleware.Auth()).Get("/youtube/videos", youtubeController.ListVideo)
		router.Middleware(middleware.Auth()).Post("/youtube/add-video", youtubeController.AddVideo)

		// Setting video
		router.Middleware(middleware.Auth()).Resource("/video-settings", settingVideosController)

		// scan komentar (mode manual)
		router.Middleware(middleware.Auth()).Get("/scan-comment/{id}", scanCommentController.ScanComment)

	})

	// Endpoint auth google
	//facades.Route().Get("/auth/google/url", )
}
