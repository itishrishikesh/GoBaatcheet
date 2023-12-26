package handlers

import "github.com/gofiber/fiber/v2"

func Welcome(context *fiber.Ctx) error {
	return context.Render("Welcome", nil, "layouts/main")
}
