package boards

import (
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
	return nil
}

func (h *Handler) Update(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	return nil
}
