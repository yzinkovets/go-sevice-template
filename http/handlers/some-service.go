package handlers

import (
	"fmt"
	"go-service-template/services"
	"go-service-template/utils/validate"

	"github.com/gofiber/fiber/v3"
)

type Request struct {
	Mac    string `json:"mac" validate:"mac,required"`
	GwUUID string `json:"gw_uuid" validate:"gw_uuid,required"`
}

type Response struct {
	Data map[string]any `json:"data"`
}

func HandleSomeServiceCall(service *services.SomeService) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Accepts(fiber.MIMEApplicationJSON)

		r := &Request{}
		if err := c.Bind().Body(r); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("invalid json body. Error: %s", err)})
		}

		if err := validate.Struct(r); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("request is not valid. Error: %s", err)})
		}

		// do some work ...
		service.Call("gw_uuid", "mac")

		ret := &Response{
			Data: map[string]any{
				"mac":    r.Mac,
				"gwUUID": r.GwUUID,
			},
		}

		return c.Status(fiber.StatusOK).JSON(ret)
	}
}
