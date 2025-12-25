package controllers

import (
	"goravel/app/http/requests"
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AuthController struct {
	// Dependent services
}

func NewAuthController() *AuthController {
	return &AuthController{
		// Inject services
	}
}

func (r *AuthController) Login(ctx http.Context) http.Response {
	// 1. Validasi Input
	var loginReq requests.LoginRequest
	errors, err := ctx.Request().ValidateRequest(&loginReq)
	if err != nil || errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, errors.All())
	}

	// 2. Cari User berdasarkan Email
	var user models.User
	err = facades.Orm().Query().Where("email", ctx.Request().Input("email")).First(&user)

	// Jika user tidak ditemukan
	if err != nil || user.ID == 0 {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{
			"message": "Kredensial tidak valid.",
		})
	}

	// 3. Verifikasi Password (Bcrypt)
	if !facades.Hash().Check(ctx.Request().Input("password"), user.Password) {
		return ctx.Response().Json(http.StatusUnauthorized, map[string]string{
			"message": "Kredensial tidak valid.",
		})
	}

	// 4. Generate Token JWT
	// Token ini akan mengikuti aturan JWT_TTL (24 jam) yang sudah kamu set
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"message": "Gagal membuat token.",
		})
	}

	// 5. Kembalikan Token dan Data User (termasuk Role)
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"token": token,
		"user": map[string]any{
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (r *AuthController) Register(ctx http.Context) http.Response {
	// 1. Validasi Input
	var registerReq requests.RegisterRequest
	errors, err := ctx.Request().ValidateRequest(&registerReq)
	if err != nil || errors != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, errors.All())
	}

	email := ctx.Request().Input("email")

	// 2. Cek apakah email sudah terdaftar (Anti-collision)
	var user models.User
	facades.Orm().Query().Where("email", email).First(&user)
	if user.ID != 0 {
		return ctx.Response().Json(http.StatusConflict, map[string]string{
			"message": "Email sudah terdaftar, silakan login.",
		})
	}

	// 3. Hash Password
	hashedPassword, err := facades.Hash().Make(ctx.Request().Input("password"))
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"message": "Gagal memproses password.",
		})
	}

	// 4. Simpan ke Database
	newUser := models.User{
		Name:     ctx.Request().Input("name"),
		Email:    email,
		Password: hashedPassword,
	}

	if err := facades.Orm().Query().Create(&newUser); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"message": "Gagal membuat akun.",
		})
	}

	// 5. Generate JWT Token langsung setelah register (UX agar user tidak perlu login lagi)
	token, err := facades.Auth(ctx).Login(&newUser)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"message": "Gagal membuat token akses.",
		})
	}

	return ctx.Response().Json(http.StatusCreated, map[string]any{
		"message": "Registrasi berhasil",
		"token":   token,
		"user":    newUser,
	})
}

func (r *AuthController) Info(ctx http.Context) http.Response {
	var user models.User

	if guard := ctx.Request().Header("Guard"); guard == "" {
		if err := facades.Auth(ctx).User(&user); err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"error": err.Error(),
			})
		}
	} else {
		if err := facades.Auth(ctx).Guard(guard).User(&user); err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"error": err.Error(),
			})
		}
	}

	return ctx.Response().Success().Json(http.Json{
		"user": user,
	})
}

func (r *AuthController) Index(ctx http.Context) http.Response {
	return nil
}
