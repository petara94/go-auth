package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	serv_dto "github.com/petara94/go-auth/internal/services/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"net/http"
)

func GetUserSelfHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			user *serv_dto.User
			err  error
		)

		session, ok := ctx.Locals(AuthSessionKey).(serv_dto.Session)
		if !ok {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(errors.New("session not found")))
		}

		user, err = userService.GetByID(session.UserID)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(user)
	}
}

func UserSelfChangePasswordHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err error
			req = &dto.ChangePassword{}
		)

		session, ok := ctx.Locals(AuthSessionKey).(serv_dto.Session)
		if !ok {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(errors.New("session not found")))
		}

		err = ctx.BodyParser(req)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		err = userService.UpdatePassword(session.UserID, req.OldPassword, req.NewPassword)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(dto.SuccessMessage())
	}
}
