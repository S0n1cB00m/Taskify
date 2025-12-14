package app

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"Taskify/internal/boards"
	"Taskify/internal/columns"
	"Taskify/internal/tasks"
	"Taskify/internal/users"
)

// Run - единственная публичная функция пакета.
// Она инициализирует зависимости и запускает приложение.
func Run() error {
	// 1. Конфигурация (хардкод или чтение из env)
	const dbDSN = "postgres://user:password@localhost:5432/dbname"
	const serverPort = ":3000"

	// 1. Настройка формата (JSON для прода, Console для локальной разработки)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Глобальная настройка уровня логирования
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Инициализация самого логгера
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// 2. Устанавливаем его как глобальный дефолтный (опционально, но удобно)
	log.Logger = logger

	logger.Info().Msg("Application starting...")

	ctx := context.Background()

	// 2. Подключение к БД
	pool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer pool.Close() // Закроем пул, когда fiber перестанет слушать порт

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	// 3. Сборка зависимостей (DI)
	userRepo := users.NewRepository(pool)
	userUseCase := users.NewUseCase(userRepo)
	userHandler := users.NewHandler(userUseCase)

	boardRepo := boards.NewRepository(pool)
	boardUseCase := boards.NewUseCase(boardRepo)
	boardHandler := boards.NewHandler(boardUseCase)

	columnRepo := columns.NewRepository(pool)
	columnUseCase := columns.NewUseCase(columnRepo)
	columnHandler := columns.NewHandler(columnUseCase)

	taskRepo := tasks.NewRepository(pool)
	taskUseCase := tasks.NewUseCase(taskRepo)
	taskHandler := tasks.NewHandler(taskUseCase)

	// 4. Инициализация Fiber
	app := fiber.New()

	// Регистрация роутов
	api := app.Group("/api")

	userHandler.RegisterRoutes(api)
	boardHandler.RegisterRoutes(api)
	columnHandler.RegisterRoutes(api)
	taskHandler.RegisterRoutes(api)

	// 5. Запуск сервера
	log.Printf("App is running on %s", serverPort)
	return app.Listen(serverPort)
}
