package app

import (
	"Taskify/internal/config"
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
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config init failed: %w", err)
	}

	logger := initLogger()
	logger.Info().Msg("Initializing application...")

	ctx := context.Background()
	pool, err := initDatabase(ctx, cfg.ConnectionURL(), logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer pool.Close()

	app := createFiberApp()

	registerRoutes(app, pool)

	logger.Info().Msgf("Server starting on :%s", cfg.HTTP.Port)
	logger.Info().Msgf("Swagger UI: http://%s:%s/swagger/index.html", cfg.HTTP.Host, cfg.HTTP.Port)

	return app.Listen(":3000")
}

func initLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()
	log.Logger = logger
	return logger
}

func initDatabase(ctx context.Context, connString string, logger zerolog.Logger) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to ping database")
		return nil, err
	}

	logger.Info().Msg("Database connection established")
	return pool, nil
}

func createFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
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

	return app
}

func registerRoutes(app *fiber.App, pool *pgxpool.Pool) {
	api := app.Group("/api")

	registerUsersRoutes(api, pool)
	registerBoardsRoutes(api, pool)
	registerColumnsRoutes(api, pool)
	registerTasksRoutes(api, pool)
}

func registerUsersRoutes(api fiber.Router, pool *pgxpool.Pool) {
	repo := users.NewRepository(pool)
	useCase := users.NewUseCase(repo)
	handler := users.NewHandler(useCase)
	handler.RegisterRoutes(api)
}

func registerBoardsRoutes(api fiber.Router, pool *pgxpool.Pool) {
	repo := boards.NewRepository(pool)
	useCase := boards.NewUseCase(repo)
	handler := boards.NewHandler(useCase)
	handler.RegisterRoutes(api)
}

func registerColumnsRoutes(api fiber.Router, pool *pgxpool.Pool) {
	repo := columns.NewRepository(pool)
	useCase := columns.NewUseCase(repo)
	handler := columns.NewHandler(useCase)
	handler.RegisterRoutes(api)
}

func registerTasksRoutes(api fiber.Router, pool *pgxpool.Pool) {
	repo := tasks.NewRepository(pool)
	useCase := tasks.NewUseCase(repo)
	handler := tasks.NewHandler(useCase)
	handler.RegisterRoutes(api)
}
