package file

import (
	"github.com/gofiber/fiber/v2"
)

func (h *FileHandler) GetFiles(c *fiber.Ctx) error {
	files, err := h.fileService.GetFiles(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get files",
		})
	}

	return c.JSON(fiber.Map{
		"files": files,
	})
}
