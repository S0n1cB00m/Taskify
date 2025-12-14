package users

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	useCase UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{
		useCase: uc,
	}
}

func (h *Handler) RegisterRoutes(app fiber.Router) {
	users := app.Group("/users")

	users.Post("/", h.Create)
	users.Get("/:id", h.GetByID)
	users.Put("/:id", h.Update)
	users.Delete("/:id", h.Delete)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	user, err := h.useCase.GetByID(c.UserContext(), userID)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.JSON(user)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	return nil
}
