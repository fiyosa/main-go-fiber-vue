package auth_repository

import (
	"go-fiber-svelte/internal/config"
	"go-fiber-svelte/internal/db"
	"go-fiber-svelte/internal/db/models"
	"go-fiber-svelte/internal/helper"
	"go-fiber-svelte/internal/http/request/auth_request"
	"go-fiber-svelte/internal/lang"
	"go-fiber-svelte/internal/lib"

	"github.com/gofiber/fiber/v2"
)

func LoginRepository(c *fiber.Ctx) error {
	req := new(auth_request.LoginRequest)

	if err := lib.Validate.Check(c, req); err != nil {
		return err
	}

	var user models.User
	result := db.RUN.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(helper.Res.Error(lang.T.Get().AUTH_FAILED, nil))
	}
	if !lib.Hash.Verify(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(helper.Res.Error(lang.T.Get().AUTH_FAILED, nil))
	}
	token, err := lib.Jwt.Create(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.Res.Error(lang.T.Get().SOMETHING_WENT_WRONG, nil))
	}
	authRecord := models.Auth{
		UserID:    user.ID,
		Token:     token,
		Revoke:    false,
		IP:        c.IP(),
		UserAgent: c.Get("User-Agent"),
	}
	db.RUN.Create(&authRecord)

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   config.APP_ENV != "local",
	})
	return c.JSON(helper.Res.Success(lang.T.Convert(lang.T.Get().SAVED_SUCCESSFULLY, map[string]any{"operator": "Login"})))
}
