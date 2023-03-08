package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/pkg"
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
			return RestErrorFromError(err)
		}

		user, err = userService.Create(*user)
		if err != nil {
			return RestErrorFromError(err)
		}

		data, err := json.Marshal(user)
		if err != nil {
			return RestErrorFromError(err)
		}

		pkg.SendJson(ctx, data)

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
			ctx.Response().SetStatusCode(fiber.StatusBadRequest)
			return RestErrorFromError(err)
		}

		page, err = strconv.Atoi(ctx.Query(PageKey))
		if err != nil {
			ctx.Response().SetStatusCode(fiber.StatusBadRequest)
			return RestErrorFromError(err)
		}

		if perPage < 0 {
			return RestErrorf("param `%s`: %w", PerPageKey, ErrBadParam)
		}

		if page < 0 {
			return RestErrorf("param `%s`: %w", PageKey, ErrBadParam)
		}

		users, err = userService.Get(perPage, page)
		if err != nil {
			ctx.Response().SetStatusCode(fiber.StatusNotFound)
			return RestErrorf("getting users: %w", err)
		}

		data, err := json.Marshal(users)
		if err != nil {
			return RestErrorFromError(err)
		}

		pkg.SendJson(ctx, data)

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
			ctx.Response().SetStatusCode(fiber.StatusBadRequest)
			return RestErrorFromError(err)
		}

		user, err = userService.GetByID(id)
		if err != nil {
			ctx.Response().SetStatusCode(fiber.StatusNotFound)
			return RestErrorFromError(err)
		}

		data, err := json.Marshal(user)
		if err != nil {
			return RestErrorFromError(err)
		}

		pkg.SendJson(ctx, data)

		return nil
	}
}
