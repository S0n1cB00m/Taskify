package v1

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	// Импорт твоих юзкейсов и домена
	domain "Taskify/services/board-service/internal/domain/board"
	"Taskify/services/board-service/internal/usecase/board"
)

type BoardHandler struct {
	createUC *board.CreateBoardUseCase
	getUC    *board.GetBoardUseCase
	listUC   *board.ListBoardsUseCase
	updateUC *board.UpdateBoardUseCase
	deleteUC *board.DeleteBoardUseCase
}

func NewBoardHandler(api fiber.Router, createUC *board.CreateBoardUseCase, getUC *board.GetBoardUseCase, listUC *board.ListBoardsUseCase, updateUC *board.UpdateBoardUseCase, deleteUC *board.DeleteBoardUseCase) {
	handler := &BoardHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}

	// Регистрируем маршруты
	boards := api.Group("/boards")
	boards.Post("/", handler.createBoard)
	boards.Get("/:id", handler.getBoard)
	boards.Get("/", handler.listBoards)
	boards.Patch("/:id", handler.updateBoard)
	boards.Delete("/:id", handler.deleteBoard)
}

// @Summary Create a new board
// @Description Create a new board with title and description
// @Tags boards
// @Accept json
// @Produce json
// @Param request body CreateBoardRequest true "Board creation info"
// @Success 201 {object} board.Board
// @Failure 400 {object} map[string]string
// @Router /boards [post]
func (h *BoardHandler) createBoard(c *fiber.Ctx) error {
	var req CreateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Вызываем ТОТ ЖЕ usecase, что и gRPC!
	cmd := board.CreateBoardCommand{
		Title:       req.Title,
		Description: req.Description,
		OwnerID:     req.Owner,
	}

	b, err := h.createUC.Handle(c.Context(), cmd)
	if err != nil {
		// Маппинг ошибок (можно вынести в middleware)
		if errors.Is(err, domain.ErrTitleRequired) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(b)
}

// @Summary Get a board
// @Description Get a board by ID
// @Tags boards
// @Accept json
// @Produce json
// @Param id path int true "Board ID"
// @Success 200 {object} board.Board
// @Failure 400 {object} map[string]string
// @Failure 404 {object} ErrBoardNotFoundResponse
// @Failure 500 {object} map[string]string
// @Router /boards/{id} [get]
func (h *BoardHandler) getBoard(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	b, err := h.getUC.Handle(c.Context(), int64(id))
	if err != nil {
		if errors.Is(err, domain.ErrBoardNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "board not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(b)
}

// @Summary List all boards
// @Description Get a list of all boards
// @Tags boards
// @Accept json
// @Produce json
// @Success 200 {array} board.Board
// @Failure 500 {object} map[string]string
// @Router /boards [get]
func (h *BoardHandler) listBoards(c *fiber.Ctx) error {
	// 1. Вызываем UseCase
	boards, err := h.listUC.Handle(c.Context())
	if err != nil {
		// Здесь ErrBoardNotFound не ожидается (пустой список - это ок),
		// поэтому сразу 500, если что-то упало в базе
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. Отдаем JSON (Fiber сам сделает маршалинг)
	// Если boards == nil, вернется null. Если хочешь [], инициализируй слайс в репо или юзкейсе.
	return c.JSON(boards)
}

// @Summary Update a  board
// @Description Update a  board with optional fields: title, description
// @Tags boards
// @Accept json
// @Produce json
// @Param id path int true "Board ID"
// @Param request body UpdateBoardRequest true "Board update info"
// @Success 200 {object} board.Board
// @Failure 400 {object} map[string]string
// @Failure 404 {object} ErrBoardNotFoundResponse
// @Router /boards/{id} [patch]
func (h *BoardHandler) updateBoard(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req UpdateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	cmd := board.UpdateBoardCommand{
		ID:          int64(id),
		Title:       req.Title,
		Description: req.Description,
	}

	b, err := h.updateUC.Handle(c.Context(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrBoardNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(b)
}

// @Summary Delete a board
// @Description Delete a board by ID
// @Tags boards
// @Accept json
// @Produce json
// @Param id path int true "Board ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /boards/{id} [delete]
func (h *BoardHandler) deleteBoard(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	err = h.deleteUC.Handle(c.Context(), int64(id))
	if err != nil {
		if errors.Is(err, domain.ErrBoardNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "board not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
