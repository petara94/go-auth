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

func SetAdminUserByIDHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err     error
			id      uint64
			isAdmin bool
		)

		admin, ok := ctx.Locals(UserAdminKey).(serv_dto.User)
		if !ok {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("user not found"))
		}

		id, err = pkg.ParseUInt64(ctx.Params("id"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		isAdmin, err = strconv.ParseBool(ctx.Query("admin"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		if !isAdmin {
			if admin.Login == AdminLoginDefault && id == admin.ID {
				return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("`admin` user can't be changed"))
			}

			currerntUser, err := userService.GetByID(id)
			if err != nil {
				return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
			}

			if admin.Login != AdminLoginDefault && id != admin.ID && currerntUser.IsAdmin {
				return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("only can change only yourself"))
			}
		}

		err = userService.SetAdmin(id, isAdmin)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.SendStatus(http.StatusOK)
	}
}

func SetBlockUserByIDHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err     error
			id      uint64
			isBlock bool
		)

		admin, ok := ctx.Locals(UserAdminKey).(serv_dto.User)
		if !ok {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("user not found"))
		}

		id, err = pkg.ParseUInt64(ctx.Params("id"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		isBlock, err = strconv.ParseBool(ctx.Query("block"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		if admin.Login == AdminLoginDefault && id == admin.ID {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("`admin` user can't be changed"))
		}

		currerntUser, err := userService.GetByID(id)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		if admin.Login != AdminLoginDefault && id != admin.ID && currerntUser.IsAdmin {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("only can change only yourself"))
		}

		err = userService.SetBlockUser(id, isBlock)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.SendStatus(http.StatusOK)
	}
}

func SetCheckPasswordUserByIDHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err     error
			id      uint64
			isCheck bool
		)

		admin, ok := ctx.Locals(UserAdminKey).(serv_dto.User)
		if !ok {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("user not found"))
		}

		id, err = pkg.ParseUInt64(ctx.Params("id"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		isCheck, err = strconv.ParseBool(ctx.Query("check"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		currerntUser, err := userService.GetByID(id)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		if id != admin.ID && currerntUser.IsAdmin {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("only can change only yourself"))
		}

		err = userService.SetCheckPassword(id, isCheck)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.SendStatus(http.StatusOK)
	}
}

func CreateEmptyUserHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err error
		)

		login := ctx.Query("login")
		if login == "" {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("login is empty"))
		}

		user, err := userService.CreateWithLogin(login)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(user)
	}
}

func GetAllUsersHandler(userService UserService) fiber.Handler {
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

func GetAllSessionsHandler(authService AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			sessions []*serv_dto.Session
			err      error
			perPage  int
			page     int
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

		sessions, err = authService.GetWithPagination(perPage, page)
		if err != nil {
			return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(sessions)
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

func UpdateUserByIDHandler(userService UserService, authService AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			putUser = &serv_dto.User{}
			err     error
			id      uint64
		)

		admin, ok := ctx.Locals(UserAdminKey).(serv_dto.User)
		if !ok {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("putUser not found"))
		}

		id, err = pkg.ParseUInt64(ctx.Params("id"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		err = ctx.BodyParser(putUser)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}
		putUser.ID = id

		// admin can't change login, isBlocked, isAdmin for himself
		if admin.Login == AdminLoginDefault && admin.ID == putUser.ID {
			putUser.Login = AdminLoginDefault
			putUser.IsAdmin = true
			putUser.Login = AdminLoginDefault
			putUser.IsBlocked = false
		}

		if admin.Login != AdminLoginDefault && admin.ID == putUser.ID {
			putUser.IsAdmin = true
			putUser.IsBlocked = false
		}

		if admin.Login != AdminLoginDefault && admin.ID != putUser.ID {
			currentUser, err := userService.GetByID(putUser.ID)
			if err != nil {
				return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
			}

			if currentUser.IsAdmin {
				return ctx.Status(http.StatusForbidden).JSON(RestErrorf("only `admin` can change admins"))
			}
		}

		// if blocked, then logout
		if putUser.IsBlocked {
			err = authService.DeleteByUserID(putUser.ID)
			if err != nil {
				return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
			}
		}

		putUser, err = userService.Update(*putUser)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return ctx.Status(http.StatusNotFound).JSON(RestErrorFromError(err))
			}
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(putUser)
	}
}

func DeleteUserByIDHandler(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var (
			err error
			id  uint64
		)

		admin, ok := ctx.Locals("user").(serv_dto.User)
		if !ok {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorf("user not found"))
		}

		id, err = pkg.ParseUInt64(ctx.Params("id"))
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		if admin.Login == AdminLoginDefault || admin.ID == id {
			return ctx.Status(http.StatusForbidden).JSON(RestErrorf("admin can't delete himself"))
		}

		user, err := userService.GetByID(id)
		if err != nil {
			return ctx.Status(http.StatusOK).JSON(dto.SuccessMessage())
		}

		// another admin can't delete main admin
		if user.Login == AdminLoginDefault {
			return ctx.Status(http.StatusForbidden).JSON(RestErrorf("main admin can't be deleted"))
		}

		err = userService.Delete(id)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(RestErrorFromError(err))
		}

		return ctx.Status(http.StatusOK).JSON(dto.SuccessMessage())
	}
}
