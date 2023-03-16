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

func CreateUserGroupHandler(userGroupService UserGroupService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			userGroup = &dto.UserGroup{}
			err       error
		)

		err = ctx.BodyParser(userGroup)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		userGroup, err = userGroupService.Create(*userGroup)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(userGroup)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJsonb(ctx, data)

		return nil
	}
}

func GetUserGroupAllHandler(userGroupService UserGroupService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			userGroups []*dto.UserGroup
			err        error
			perPage    int
			page       int
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

		userGroups, err = userGroupService.Get(perPage, page)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(userGroups)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJsonb(ctx, data)

		return nil
	}
}

func GetUserGroupByIDHandler(userGroupService UserGroupService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			userGroup *dto.UserGroup
			err       error
			id        uint64
		)

		id, err = strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		userGroup, err = userGroupService.GetByID(id)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(userGroup)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJsonb(ctx, data)

		return nil
	}
}

func UpdateUserGroupHandler(userGroupService UserGroupService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			userGroup = &dto.UserGroup{}
			err       error
			id        uint64
		)

		id, err = strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		err = ctx.BodyParser(userGroup)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		userGroup.ID = id

		userGroup, err = userGroupService.Update(*userGroup)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
			}
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		data, err := json.Marshal(userGroup)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(RestErrorFromError(err))
		}

		SendJsonb(ctx, data)

		return nil
	}
}

func DeleteUserGroupHandler(userGroupService UserGroupService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err error
			id  uint64
		)

		id, err = strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		err = userGroupService.Delete(id)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(dto.SuccessMessage())
	}
}
