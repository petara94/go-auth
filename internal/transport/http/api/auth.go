package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	serv_dto "github.com/petara94/go-auth/internal/services/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/mappers"
	"net/http"
)

func LoginHandler(authService AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			auth = &dto.LoginReq{}
			err  error
		)

		err = ctx.BodyParser(auth)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		session, err := authService.Login(*mappers.LoginReqToAuth(auth))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(session)
	}
}

func RegisterHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			auth = &serv_dto.Auth{}
			err  error
		)

		err = ctx.BodyParser(auth)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		_, err = userService.CreateWithLoginAndPassword(auth.Login, auth.Password)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.SendStatus(http.StatusOK)
	}
}

func LogoutHandler(userService AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err error
		)

		session, ok := ctx.Locals(AuthSessionKey).(serv_dto.Session)
		if !ok {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(errors.New("session not found")))
		}

		err = userService.Logout(serv_dto.Session(session))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.SendStatus(http.StatusOK)
	}
}
