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

// Create
// @Summary      Создание пользователя
// @Description  Создает пользователя, используя поля: email, username, password
// @Tags         users
// @Accept		 json
// @Produce      json
// @Param        RequestBody body CreateUserDTO true "Данные для регистрации"
// @Success      201  {object}  User
// @Router       /users [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var dto CreateUserDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	user := &User{
		Email:    dto.Email,
		Username: dto.Username,
		Password: dto.Password,
	}

	createdUser, err := h.useCase.Create(c.UserContext(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

// GetByID
// @Summary      Получение пользователя
// @Description  Получает пользователя по ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "id пользователя"
// @Success      200  {object}  User
// @Router       /users/{id} [get]
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

	return c.Status(fiber.StatusOK).JSON(user)
}

// Update
// @Summary      Обновление пользователя
// @Description  Обновления пользовательских полей: email, username, password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      int                   true  "id пользователя"
// @Param        input body      CreateUserDTO true  "Поля для редактирования"
// @Success      200   {object}  User
// @Router       /users/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	var dto CreateUserDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	user := &User{
		ID:       userID,
		Email:    dto.Email,
		Username: dto.Username,
		Password: dto.Password,
	}

	createdUser, err := h.useCase.Update(c.UserContext(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(createdUser)
}

// Delete
// @Summary      Удаление пользователя
// @Description  Удаляет пользователя по ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  User
// @Router       /users/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(user)
}
