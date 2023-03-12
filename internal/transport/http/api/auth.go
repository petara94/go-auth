package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"net/http"
)

func LoginHandler(authService AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			auth = &dto.Auth{}
			err  error
		)

		err = ctx.BodyParser(auth)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		session, err := authService.Login(*auth)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(session)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJsonb(ctx, data)

		return nil
	}
}

func LogoutHandler(userService AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			auth = &dto.Auth{}
			err  error
		)

		err = ctx.BodyParser(auth)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		session, err := userService.Login(*auth)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(session)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJsonb(ctx, data)

		return nil
	}
}
