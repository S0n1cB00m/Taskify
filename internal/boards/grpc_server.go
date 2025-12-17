package boards

import (
	"context"
	"errors"

	boardspb "Taskify/internal/pb/boards"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	boardspb.UnimplementedBoardsServiceServer
	uc UseCase
}

func NewGRPCServer(uc UseCase) *GRPCServer {
	return &GRPCServer{uc: uc}
}

// CreateBoard
func (s *GRPCServer) CreateBoard(ctx context.Context, req *boardspb.CreateBoardRequest) (*boardspb.CreateBoardResponse, error) {
	log.Ctx(ctx).Info().
		Int64("user_id", req.UserId).
		Msg("boards-service: CreateBoard called")

	board := &Board{
		Name:        req.Name,
		Description: req.Description,
		UserId:      req.UserId,
	}

	created, err := s.uc.Create(ctx, board)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("boards-service: Create failed")
		return nil, status.Errorf(codes.Internal, "create board: %v", err)
	}

	return &boardspb.CreateBoardResponse{
		Board: &boardspb.Board{
			Id:          created.Index, // или created.ID, зависит от твоей модели
			Name:        created.Name,
			Description: created.Description,
		},
	}, nil
}

// GetBoardByID
func (s *GRPCServer) GetBoardByID(ctx context.Context, req *boardspb.GetBoardByIDRequest) (*boardspb.GetBoardResponse, error) {
	log.Ctx(ctx).Info().
		Int64("user_id", req.UserId).
		Int64("id", req.Id).
		Msg("boards-service: GetBoardByID called")

	board := &Board{
		Index:  req.Id,
		UserId: req.UserId,
	}

	got, err := s.uc.GetByID(ctx, board)
	if err != nil {
		if errors.Is(err, ErrBoardNotFound) {
			return nil, status.Errorf(codes.NotFound, "board not found")
		}
		log.Ctx(ctx).Error().Err(err).Msg("boards-service: GetByID failed")
		return nil, status.Errorf(codes.Internal, "get board: %v", err)
	}

	return &boardspb.GetBoardResponse{
		Board: &boardspb.Board{
			Id:          got.Index,
			Name:        got.Name,
			Description: got.Description,
		},
	}, nil
}

// UpdateBoard
func (s *GRPCServer) UpdateBoard(ctx context.Context, req *boardspb.UpdateBoardRequest) (*boardspb.UpdateBoardResponse, error) {
	log.Ctx(ctx).Info().
		Int64("user_id", req.UserId).
		Int64("id", req.Id).
		Msg("boards-service: UpdateBoard called")

	board := &Board{
		Index:       req.Id,
		Name:        req.Name,
		Description: req.Description,
		UserId:      req.UserId,
	}

	updated, err := s.uc.Update(ctx, board)
	if err != nil {
		if errors.Is(err, ErrBoardNotFound) {
			return nil, status.Errorf(codes.NotFound, "board not found")
		}
		log.Ctx(ctx).Error().Err(err).Msg("boards-service: Update failed")
		return nil, status.Errorf(codes.Internal, "update board: %v", err)
	}

	return &boardspb.UpdateBoardResponse{
		Board: &boardspb.Board{
			Id:          updated.Index,
			Name:        updated.Name,
			Description: updated.Description,
		},
	}, nil
}

func (s *GRPCServer) DeleteBoard(ctx context.Context, req *boardspb.DeleteBoardRequest) (*emptypb.Empty, error) {
	log.Ctx(ctx).Info().
		Int64("user_id", req.UserId).
		Int64("id", req.Id).
		Msg("boards-service: DeleteBoard called")

	board := &Board{
		Index:  req.Id,
		UserId: req.UserId,
	}

	if err := s.uc.Delete(ctx, board); err != nil {
		if errors.Is(err, ErrBoardNotFound) {
			return nil, status.Errorf(codes.NotFound, "board not found")
		}
		log.Ctx(ctx).Error().Err(err).Msg("boards-service: Delete failed")
		return nil, status.Errorf(codes.Internal, "delete board: %v", err)
	}

	return &emptypb.Empty{}, nil
}
