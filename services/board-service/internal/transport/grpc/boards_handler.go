package grpc_handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	// 1. Импортируем сгенерированный код Protobuf (Алиас pb для краткости)
	// Замени путь на свой реальный путь к сгенерированным файлам
	pb "Taskify/proto/boards/v1"

	// 2. Импортируем UseCase (Бизнес-сценарии)
	usecase "Taskify/services/board-service/internal/usecase/board"

	// 3. Импортируем Домен (для проверки типов ошибок)
	domain "Taskify/services/board-service/internal/domain/board"
)

// Handler реализует интерфейс, который сгенерировал protoc (UnimplementedBoardServiceServer)
type Handler struct {
	// Встраиваем структуру для совместимости (обязательно для gRPC Go)
	pb.UnimplementedBoardServiceServer

	// Зависимость: Хендлер знает только про UseCase
	createBoardUC *usecase.CreateBoardUseCase
	getBoardUC    *usecase.GetBoardUseCase
	listBoardsUC  *usecase.ListBoardsUseCase
	updateBoardUC *usecase.UpdateBoardUseCase
	deleteBoardUC *usecase.DeleteBoardUseCase
}

// Конструктор
func NewHandler(createUC *usecase.CreateBoardUseCase, getUC *usecase.GetBoardUseCase, listBoardsUC *usecase.ListBoardsUseCase, updateBoardUC *usecase.UpdateBoardUseCase, deleteUC *usecase.DeleteBoardUseCase) *Handler {
	return &Handler{
		createBoardUC: createUC,
		getBoardUC:    getUC,
		listBoardsUC:  listBoardsUC,
		updateBoardUC: updateBoardUC,
		deleteBoardUC: deleteUC,
	}
}

func toProtoBoard(b *domain.Board) *pb.Board {
	return &pb.Board{
		Id:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Owner:       b.Owner,
		CreatedAt:   timestamppb.New(b.CreatedAt),
		UpdatedAt:   timestamppb.New(b.UpdatedAt),
	}
}

// CreateBoard — это метод, который вызовет gRPC сервер, когда придет запрос
func (h *Handler) CreateBoard(ctx context.Context, req *pb.CreateBoardRequest) (*pb.CreateBoardResponse, error) {
	// ШАГ 1: Преобразуем gRPC Request -> UseCase Command
	command := usecase.CreateBoardCommand{
		Title:       req.Title,
		Description: req.Description,
		OwnerID:     req.Owner,
	}

	// ШАГ 2: Вызываем бизнес-логику
	createdBoard, err := h.createBoardUC.Handle(ctx, command)

	// ШАГ 3: Обрабатываем ошибки
	if err != nil {
		// Пытаемся понять, какая именно ошибка произошла, чтобы вернуть верный HTTP/gRPC код
		switch {
		case errors.Is(err, domain.ErrTitleRequired), errors.Is(err, domain.ErrTitleTooLong), errors.Is(err, domain.ErrEmptyOwner):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "internal error: %v", err)
		}
	}

	var board pb.Board

	board.Id = createdBoard.ID
	board.Title = createdBoard.Title
	board.Description = createdBoard.Description
	board.Owner = createdBoard.Owner
	board.CreatedAt = timestamppb.New(createdBoard.CreatedAt)
	board.UpdatedAt = timestamppb.New(createdBoard.UpdatedAt)

	// ШАГ 4: Преобразуем Доменную сущность -> gRPC Response
	return &pb.CreateBoardResponse{
		Board: toProtoBoard(createdBoard),
	}, nil
}

func (h *Handler) GetBoard(ctx context.Context, id *pb.GetBoardRequest) (*pb.GetBoardResponse, error) {
	receivedBoard, err := h.getBoardUC.Handle(ctx, id.Id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBoardNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "internal error: %v", err)
		}
	}

	return &pb.GetBoardResponse{
		Board: toProtoBoard(receivedBoard),
	}, nil
}

func (h *Handler) ListBoards(ctx context.Context, _ *emptypb.Empty) (*pb.ListBoardsResponse, error) {
	// 1. Получаем доменные сущности
	domainBoards, err := h.listBoardsUC.Handle(ctx)
	if err != nil {
		// ErrBoardNotFound тут вряд ли будет (пустой список != ошибка), но оставим для примера
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	// 2. Конвертируем []*domain.Board -> []*pb.Board
	protoBoards := make([]*pb.Board, 0, len(domainBoards))
	for _, b := range domainBoards {
		protoBoards = append(protoBoards, toProtoBoard(b))
	}

	// 3. Возвращаем результат
	return &pb.ListBoardsResponse{
		Board: protoBoards, // Поле в proto называется 'repeated Board board = 1;'
	}, nil
}

func (h *Handler) UpdateBoard(ctx context.Context, req *pb.UpdateBoardRequest) (*pb.UpdateBoardResponse, error) {
	cmd := usecase.UpdateBoardCommand{
		ID:          req.Id,
		Title:       req.Title,       // Это уже *string благодаря 'optional' в proto
		Description: req.Description, // Это тоже *string
	}

	updatedBoard, err := h.updateBoardUC.Handle(ctx, cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &pb.UpdateBoardResponse{
		Board: toProtoBoard(updatedBoard),
	}, nil
}

func (h *Handler) DeleteBoard(ctx context.Context, id *pb.DeleteBoardRequest) (*emptypb.Empty, error) {
	err := h.deleteBoardUC.Handle(ctx, id.Id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBoardNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "internal error: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}
