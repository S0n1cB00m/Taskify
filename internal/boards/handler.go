package boards

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

	users.Post("/:user_id/boards", h.Create)
	users.Get("/:user_id/boards/:id", h.GetByID)
	users.Put("/:user_id/boards/:id", h.Update)
	users.Delete("/:user_id/boards/:id", h.Delete)
}

// Create
// @Summary      Создание таблицы
// @Description  Создает пользователя, используя id пользователя и поля: name, description
// @Tags         users
// @Accept		 json
// @Produce      json
// @Param        user_id     path      int  true  "id пользователя"
// @Param        RequestBody body CreateBoardDTO true "Данные для создания"
// @Success      201  {object}  Board
// @Router        /users/{user_id}/boards [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Params("board_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	var dto CreateBoardDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	board := &Board{
		Name:        dto.Name,
		Description: dto.Description,
		UserId:      userId,
	}

	createdBoard, err := h.useCase.Create(c.UserContext(), board)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create board",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdBoard)
}

// GetByID
// @Summary      Получение таблицы
// @Description  Получает таблицы пользователя по ее порядковому номеру
// @Tags         users
// @Produce      json
// @Param        user_id   path      int  true  "id пользователя"
// @Param        id        path      int  true  "id таблицы"
// @Success      200  {object}  Board
// @Router       /users/{user_id}/boards/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Params("board_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	boardId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	board := &Board{
		Index:  boardId,
		UserId: userId,
	}

	receivedBoard, err := h.useCase.GetByID(c.UserContext(), board)
	if err != nil {
		if errors.Is(err, ErrBoardNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Board not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(receivedBoard)
}

// Update
// @Summary      Обновление таблицы
// @Description  Обновления полей таблицы: name, description
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id   path      int  true  "id пользователя"
// @Param        id        path      int  true  "id таблицы"
// @Param        input body      CreateBoardDTO true  "Поля для редактирования"
// @Success      200   {object}  Board
// @Router       /users/{user_id}/boards/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Params("board_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	boardId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	var dto CreateBoardDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	board := &Board{
		Index:       boardId,
		Name:        dto.Name,
		Description: dto.Description,
		UserId:      userId,
	}

	updatedBoard, err := h.useCase.Update(c.UserContext(), board)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update board",
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedBoard)
}

// Delete
// @Summary      Удаление таблицы
// @Description  Удаляет таблицу пользователя по ее порядковому номеру
// @Tags         users
// @Produce      json
// @Param        user_id   path      int  true  "id пользователя"
// @Param        id        path      int  true  "id таблицы"
// @Success      200  {object}  Board
// @Router       /users/{user_id}/boards/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	boardId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	board := &Board{
		Index:  boardId,
		UserId: userId,
	}

	deletedBoard, err := h.useCase.GetByID(c.UserContext(), board)
	if err != nil {
		if errors.Is(err, ErrBoardNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Board not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(deletedBoard)
}
