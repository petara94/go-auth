package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/petara94/go-auth/assets"
	"github.com/petara94/go-auth/internal/services/dto"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func LoggerMiddleware(logger zap.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := ctx.Request()
		res := ctx.Response()
		start := time.Now()

		// Latency
		latency := time.Since(start)

		// Request id
		requestID := string(req.Header.Peek(fiber.HeaderXRequestID))

		// Request path
		requestPath := string(req.URI().Path())
		if requestPath == "" {
			requestPath = "/"
		}

		// Bytes in
		bytesIn := string(req.Header.Peek(fiber.HeaderContentLength))
		if bytesIn == "" {
			bytesIn = "0"
		}

		// Bytes out
		bytesOut := string(res.Header.Peek(fiber.HeaderContentLength))
		if bytesOut == "" {
			bytesOut = "0"
		}

		// Log
		logger.Debug(
			"request",
			zap.String("remote_ip", ctx.IP()),
			zap.Time("date", time.Now()),
			zap.String("method", ctx.Method()),
			zap.String("proto", ctx.Protocol()),
			zap.Int("status", res.StatusCode()),
			zap.String("bytes_in", bytesIn),
			zap.String("bytes_out", bytesOut),
			zap.String("latency_human", latency.String()),
			zap.String("latency_ns", strconv.FormatInt(int64(latency), 10)),
			zap.String("id", requestID),
			zap.String("host", string(req.Host())),
			zap.String("uri", string(req.RequestURI())),
			zap.String("path", requestPath),
		)
		return ctx.Next()
	}
}

func CheckAuthorizeMiddleware(authService AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		token, ok := headers[TokenHeader]

		if !ok {
			return ctx.Status(http.StatusForbidden).JSON(RestErrorFromError(ErrNotAuthorised))
		}

		session, err := authService.Get(token)
		if err != nil {
			return ctx.Status(http.StatusForbidden).JSON(RestErrorFromError(ErrTokenExpired))
		}

		if session.Expr != nil {
			if time.Now().Unix() > session.Expr.Unix() {
				return ctx.Status(http.StatusForbidden).JSON(RestErrorFromError(ErrTokenExpired))
			}
		}

		ctx.Locals(AuthSessionKey, *session)

		return ctx.Next()
	}
}

// CheckAdminMiddleware middlware check is user admin
func CheckAdminMiddleware(userService UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		session, ok := ctx.Locals(AuthSessionKey).(dto.Session)
		if !ok {
			return ctx.Status(http.StatusForbidden).JSON(RestErrorFromError(ErrNotAuthorised))
		}

		// get user from session
		user, err := userService.GetByID(session.UserID)
		if err != nil {
			return ctx.Status(http.StatusForbidden).JSON(RestErrorFromError(ErrNotAuthorised))
		}

		if !user.IsAdmin {
			return ctx.Status(http.StatusForbidden).JSON(RestErrorFromError(ErrNotAllowed))
		}

		ctx.Locals(UserAdminKey, *user)

		return ctx.Next()
	}
}

func SwaggerMiddleware() fiber.Handler {
	return filesystem.New(
		filesystem.Config{
			Root:       http.FS(assets.SwaggerFiles),
			PathPrefix: "",
			Index:      "index.html",
		},
	)
}
