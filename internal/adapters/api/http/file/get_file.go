package file

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (h *FileHandler) GetFile(c *fiber.Ctx) error {
	fileID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "incorrect id",
		})
	}
	filePath, err := h.fileService.GetFile(c.UserContext(), fileID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting file",
		})
	}

	return c.SendFile(filePath)
}
