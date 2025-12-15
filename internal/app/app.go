package app

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "Taskify/docs"

	fiberSwagger "github.com/gofiber/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"Taskify/internal/boards"
	"Taskify/internal/columns"
	"Taskify/internal/tasks"
	"Taskify/internal/users"
)

func Run() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	log.Logger = logger
	logger.Info().Msg("Initializing application...")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://postgres:AaGgRrUuNnAa99!@localhost:5433/postgres")
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer pool.Close()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Или укажите конкретно "http://localhost:3000"
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("X-Request-ID", reqID)

		l := log.With().Str("req_id", reqID).Logger()

		ctx := l.WithContext(c.UserContext())
		c.SetUserContext(ctx)

		return c.Next()
	})

	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	api := app.Group("/api")

	userRepo := users.NewRepository(pool)
	userUseCase := users.NewUseCase(userRepo)
	userHandler := users.NewHandler(userUseCase)
	userHandler.RegisterRoutes(api)

	boardRepo := boards.NewRepository(pool)
	boardUseCase := boards.NewUseCase(boardRepo)
	boardHandler := boards.NewHandler(boardUseCase)
	boardHandler.RegisterRoutes(api)

	columnRepo := columns.NewRepository(pool)
	columnUseCase := columns.NewUseCase(columnRepo)
	columnHandler := columns.NewHandler(columnUseCase)
	columnHandler.RegisterRoutes(api)

	taskRepo := tasks.NewRepository(pool)
	taskUseCase := tasks.NewUseCase(taskRepo)
	taskHandler := tasks.NewHandler(taskUseCase)
	taskHandler.RegisterRoutes(api)

	logger.Info().Msg("Server starting on :3000")
	log.Info().Msg("Swagger UI: http://185.68.22.208:3000/swagger/index.html")

	return app.Listen(":3000")
}
