package main

import(
	"github.com/gofiber/fiber/v2"
	"log"
)

func main () {
	app := fiber.New()

	app.Get("/", func (c *fiber.Ctx) error {
		return c.SendString("PropAnalytix Backend is Running")
	})

	log.Fatal(app.Listen(":3000"))
}