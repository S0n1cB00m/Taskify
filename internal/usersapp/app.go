package usersapp

import (
	"Taskify/internal/config"
	userspb "Taskify/internal/pb/users"
	"Taskify/internal/users"

	"context"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func Run() error {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("config init failed")
	}

	logger := initLogger()
	logger.Info().Msg("users-service: initializing...")

	pool, err := initDatabase(ctx, cfg.PG.ConnectionURL(), logger)
	if err != nil {
		return err
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", cfg.GRPC.UsersAddress)
	if err != nil {
		logger.Error().Err(err).Msg("listen failed")
		return err
	}

	s := grpc.NewServer()

	userRepo := users.NewRepository(pool)
	userUC := users.NewUseCase(userRepo)
	userspb.RegisterUsersServiceServer(s, users.NewGRPCServer(userUC))

	logger.Info().Msgf("users-service gRPC listening on %s", cfg.GRPC.UsersAddress)

	if err := s.Serve(lis); err != nil {
		logger.Error().Err(err).Msg("serve failed")
		return err
	}

	return nil
}

func initLogger() zerolog.Logger {
	// можно переиспользовать общий initLogger, если вынесешь его в отдельный пакет
	logger := log.Logger
	return logger
}

func initDatabase(ctx context.Context, connString string, logger zerolog.Logger) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Error().Err(err).Msg("users-service: db connect failed")
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Error().Err(err).Msg("users-service: db ping failed")
		return nil, err
	}

	logger.Info().Msg("users-service: database connection established")
	return pool, nil
}
