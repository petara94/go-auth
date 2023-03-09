package api

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/petara94/go-auth/internal/repo"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"net/http"
	"strconv"
)

func CreateUserHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			user = &dto.User{}
			err  error
		)

		err = ctx.BodyParser(user)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		user, err = userService.Create(*user)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(user)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJson(ctx, data)

		return nil
	}
}

func GetUserAllHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			users   []*dto.User
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

		users, err = userService.Get(perPage, page)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(users)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJson(ctx, data)

		return nil
	}
}

func GetUserByIDHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			user *dto.User
			err  error
			id   uint64
		)

		id, err = strconv.ParseUint(ctx.Params("id"), 10, 64)
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

		SendJson(ctx, data)

		return nil
	}
}

func UpdateUserHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			user = &dto.User{}
			err  error
			id   uint64
		)

		id, err = strconv.ParseUint(ctx.Params("id"), 10, 64)
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

		SendJson(ctx, data)

		return nil
	}
}

func DeleteUserHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err error
			id  uint64
		)

		id, err = strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		err = userService.Delete(id)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(RestErrorFromError(Success))
	}
}
