package api

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/petara94/go-auth/internal/repo"
	serv_dto "github.com/petara94/go-auth/internal/services/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/pkg"
	"net/http"
	"strconv"
)

func CreateUserHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			regUser = &serv_dto.User{}
			err     error
		)

		err = ctx.BodyParser(regUser)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		_, err = userService.Create(*regUser)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.SendStatus(http.StatusOK)
	}
}

func GetUserAllHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			users   []*serv_dto.User
			err     error
			perPage int
			page    int
		)

		perPage, err = strconv.Atoi(ctx.Query(PerPageKey))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		page, err = strconv.Atoi(ctx.Query(PageKey))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		if perPage < 0 {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("param `%s`: %s", PerPageKey, ErrBadParam))
		}

		if page < 0 {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("param `%s`: %s", PageKey, ErrBadParam))
		}

		users, err = userService.GetWithPagination(perPage, page)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(users)
	}
}

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

func GetUserByIDHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			user *serv_dto.User
			err  error
			id   uint64
		)

		id, err = pkg.ParseUInt64(ctx.Params("id"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		user, err = userService.GetByID(id)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(user)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJsonb(ctx, data)

		return nil
	}
}

func UpdateUserHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			user = &serv_dto.User{}
			err  error
			id   uint64
		)

		id, err = pkg.ParseUInt64(ctx.Params("id"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		err = ctx.BodyParser(user)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		user.ID = id

		user, err = userService.Update(*user)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
			}
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(user)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJsonb(ctx, data)

		return nil
	}
}

func DeleteUserHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err error
			id  uint64
		)

		id, err = pkg.ParseUInt64(ctx.Params("id"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		err = userService.Delete(id)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(dto.SuccessMessage())
	}
}
