package app

import (
	"Taskify/internal/columns"
	"Taskify/internal/tasks"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log" // Глобальный алиас, нужен для настройки

	// Импорты твоих сервисов
	"Taskify/internal/boards"
	"Taskify/internal/users"
)

func Run() error {
	// --- 1. Настройка ЕДИНОГО логгера ---
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Красивый вывод в консоль (для локальной разработки).
	// Для прода лучше убрать ConsoleWriter и оставить JSON.
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	// Создаем корневой логгер
	logger := zerolog.New(output).With().Timestamp().Logger()

	// Делаем его глобальным дефолтным (на всякий случай)
	log.Logger = logger

	logger.Info().Msg("Initializing application...")

	// --- 2. БД ---
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://user:password@localhost:5432/dbname")
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer pool.Close()

	// --- 3. Fiber и Middleware ---
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// ВАЖНО: Подключаем middleware, который будет "раздавать" логгер всем сервисам
	app.Use(func(c *fiber.Ctx) error {
		// Генерируем ID запроса
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("X-Request-ID", reqID)

		// Создаем child-логгер с ID запроса
		// Мы берем тот самый logger, который настроили выше (log.Logger)
		l := log.With().Str("req_id", reqID).Logger()

		// Инжектим логгер в контекст
		ctx := l.WithContext(c.UserContext())
		c.SetUserContext(ctx)

		return c.Next()
	})

	// --- 4. Сборка модулей ---
	// Обрати внимание: мы НЕ передаем logger в конструкторы!
	api := app.Group("/api")

	// --- 3. DI (Сборка слоев) ---
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
	return app.Listen(":3000")
}
