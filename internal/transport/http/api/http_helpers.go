package api

import (
	"github.com/gofiber/fiber/v2"
)

func SendJsonb(ctx *fiber.Ctx, body []byte) {
	ctx.Response().SetBody(body)
	ctx.Set("Content-Type", "application/json")
}