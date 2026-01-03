package main

//swag init -g services/board-service/cmd/main.go --output services/board-service/docs

import (
	"context"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "Taskify/services/board-service/docs"

	"github.com/gofiber/swagger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	// Импорты твоих пакетов
	// (Убедись, что пути совпадают с твоим go.mod)
	pb "Taskify/proto/boards/v1"

	"Taskify/services/board-service/internal/config"
	"Taskify/services/board-service/internal/infrastructure/persistence"
	grpcHandler "Taskify/services/board-service/internal/transport/grpc"
	httpHandler "Taskify/services/board-service/internal/transport/http/v1"
	usecaseBoard "Taskify/services/board-service/internal/usecase/board"
)

// @title           Taskify Board Service API
// @version         1.0
// @description     This is a board service for Taskify application.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            185.68.22.208:3000
// @BasePath        /v1
func main() {
	// 1. Конфигурация (пока хардкод для простоты, потом вынесем в config)
	ctx := context.Background()

	serviceConfig := config.MustLoad()

	setupLogger(serviceConfig.Env)

	// 2. Подключение к БД (PgxPool)
	poolConfig, err := pgxpool.ParseConfig(serviceConfig.Postgres.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse DB config")
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to DB: %v")
	}
	defer dbPool.Close()

	// Проверяем соединение
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("DB Ping failed: %v")
	}
	log.Printf("Successfully connected to Database")

	// 3. Инициализация слоев (Dependency Injection)

	// Layer 1: Persistence (Repository)
	boardRepo := persistence.NewBoardRepository(dbPool)

	// Layer 2: UseCase (Business Logic)
	// Тут можно создать сразу структуру, которая держит все юзкейсы,
	// но пока у нас один - инициализируем его.
	createBoardUC := usecaseBoard.NewCreateBoardUseCase(boardRepo)
	getBoardUC := usecaseBoard.NewGetBoardUseCase(boardRepo)
	listBoardsUC := usecaseBoard.NewListBoardsUseCase(boardRepo)
	updateBoardUC := usecaseBoard.NewUpdateBoardUseCase(boardRepo)
	deleteBoardUC := usecaseBoard.NewDeleteBoardUseCase(boardRepo)

	// Layer 3: Transport (gRPC Handler)
	boardHandler := grpcHandler.NewHandler(createBoardUC, getBoardUC, listBoardsUC, updateBoardUC, deleteBoardUC)

	// 4. Запуск gRPC сервера
	lis, err := net.Listen("tcp", serviceConfig.GRPC.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen: %v")
	}

	grpcServer := grpc.NewServer()

	// Регистрируем наш сервис
	pb.RegisterBoardServiceServer(grpcServer, boardHandler)

	// Включаем Reflection (чтобы можно было стучаться через Postman/gRPCui)
	reflection.Register(grpcServer)

	log.Printf("Board Service (gRPC) is running on port %s", serviceConfig.GRPC.Port)

	// ВАЖНО: Запускаем в горутине!
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC: %v")
		}
	}()

	// --- HTTP Server (Fiber) ---
	app := fiber.New()
	v1 := app.Group("/v1")

	// Теперь httpHandler будет доступен (при условии, что ты добавил import выше)
	httpHandler.NewBoardHandler(v1, createBoardUC, getBoardUC, listBoardsUC, updateBoardUC, deleteBoardUC)

	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Printf("Starting HTTP server on :3000")
	log.Err(app.Listen(serviceConfig.HTTP.Port)) // Этот вызов тоже блокирующий, так что main не завершится
}

func setupLogger(env string) {
	switch env {
	case "local":
		// Человекочитаемый вывод в консоль
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "dev":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		// По умолчанию JSON формат, удобный для Kibana/Datadog
	case "prod":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
