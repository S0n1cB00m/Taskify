package users

import (
	"errors"
	"strconv"

	userspb "Taskify/internal/pb/users"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	client userspb.UsersServiceClient
}

func NewHandler(client userspb.UsersServiceClient) *Handler {
	return &Handler{client: client}
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
	log.Ctx(c.UserContext()).Info().
		Msg("users-service: handler.Create called")

	var dto CreateUserDTO

	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	resp, err := h.client.CreateUser(
		c.UserContext(),
		&userspb.CreateUserRequest{
			Email:    dto.Email,
			Username: dto.Username,
			Password: dto.Password,
		},
	)
	if err != nil {
		log.Ctx(c.UserContext()).
			Error().
			Err(err).
			Msg("gateway: CreateUser gRPC call failed")

		if errors.Is(err, ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "user with this email already exists",
			})
		}

		st, ok := status.FromError(err)
		if ok {
			log.Ctx(c.UserContext()).
				Error().
				Str("code", st.Code().String()).
				Str("message", st.Message()).
				Msg("gateway: gRPC status")
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	log.Ctx(c.UserContext()).Info().
		Msg("users-service: handler.Create finished")

	return c.Status(fiber.StatusCreated).JSON(resp.User)
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
	userId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	resp, err := h.client.GetUserByID(
		c.UserContext(),
		&userspb.GetUserByIDRequest{Id: userId},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(resp.User)
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
	userId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	var dto CreateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	resp, err := h.client.UpdateUser(
		c.UserContext(),
		&userspb.UpdateUserRequest{
			Id:       userId,
			Email:    dto.Email,
			Username: dto.Username,
			Password: dto.Password,
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp.User)
}

// Delete
// @Summary      Удаление пользователя
// @Description  Удаляет пользователя по ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "id пользователя"
// @Success      200  {object}  User
// @Router       /users/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad request"})
	}

	_, err = h.client.DeleteUser(
		c.UserContext(),
		&userspb.DeleteUserRequest{Id: userId})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "SUCCESS"})
}
