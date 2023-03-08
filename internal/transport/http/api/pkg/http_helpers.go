package pkg

import "github.com/gofiber/fiber/v2"

func SendJson(ctx *fiber.Ctx, body []byte) {
	ctx.Response().SetBody(body)
	ctx.Set("Content-Type", "application/json")
}
