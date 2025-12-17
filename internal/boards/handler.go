package boards

import (
	"strconv"

	boardspb "Taskify/internal/pb/boards"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	client boardspb.BoardsServiceClient
}

func NewHandler(client boardspb.BoardsServiceClient) *Handler {
	return &Handler{client: client}
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
	userID, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	var dto CreateBoardDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
	}

	log.Ctx(c.UserContext()).
		Info().
		Int64("user_id", userID).
		Msg("gateway: CreateBoard called")

	resp, err := h.client.CreateBoard(
		c.UserContext(),
		&boardspb.CreateBoardRequest{
			Name:        dto.Name,
			Description: dto.Description,
			UserId:      userID,
		},
	)
	if err != nil {
		log.Ctx(c.UserContext()).
			Error().
			Err(err).
			Msg("gateway: CreateBoard gRPC failed")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create board",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp.Board)
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
	userID, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	boardIndex, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	resp, err := h.client.GetBoardByID(
		c.UserContext(),
		&boardspb.GetBoardByIDRequest{
			Id:     boardIndex,
			UserId: userID,
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Board not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(resp.Board)
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
	userID, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	boardIndex, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	var dto CreateBoardDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
	}

	resp, err := h.client.UpdateBoard(
		c.UserContext(),
		&boardspb.UpdateBoardRequest{
			Id:          boardIndex,
			Name:        dto.Name,
			Description: dto.Description,
			UserId:      userID,
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update board"})
	}

	return c.Status(fiber.StatusOK).JSON(resp.Board)
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
	userID, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	boardIndex, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	_, err = h.client.DeleteBoard(
		c.UserContext(),
		&boardspb.DeleteBoardRequest{
			Id:     boardIndex,
			UserId: userID,
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Board not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "SUCCESS"})
}
