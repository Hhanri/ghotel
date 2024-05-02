package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hhanri/ghotel/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		ID:        "ID",
		FirstName: "First",
		LastName:  "Last",
	}
	return c.JSON(user)
}

func HandleGetUser(c *fiber.Ctx) error {
	user := types.User{
		ID:        "ID",
		FirstName: "First",
		LastName:  "Last",
	}
	return c.JSON(user)
}
