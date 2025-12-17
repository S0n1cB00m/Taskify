package users

import (
	"context"
	"errors"

	userspb "Taskify/internal/pb/users"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	userspb.UnimplementedUsersServiceServer
	uc UseCase
}

func NewGRPCServer(uc UseCase) *GRPCServer {
	return &GRPCServer{uc: uc}
}

// CreateUser
func (s *GRPCServer) CreateUser(ctx context.Context, req *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {
	log.Ctx(ctx).Info().
		Str("email", req.Email).
		Msg("users-service: CreateUser called")

	user := &User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password, // usecase сам захеширует
	}

	created, err := s.uc.Create(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().
			Err(err).
			Str("email", req.Email).
			Msg("users-service: failed to create user")

		// можно делать разную мапу ошибок, пока общий Internal
		return nil, status.Errorf(codes.Internal, "create user: %v", err)
	}

	return &userspb.CreateUserResponse{
		User: &userspb.User{
			Id:       created.Id,
			Email:    created.Email,
			Username: created.Username,
		},
	}, nil
}

// GetUserByID
func (s *GRPCServer) GetUserByID(ctx context.Context, req *userspb.GetUserByIDRequest) (*userspb.GetUserResponse, error) {
	log.Ctx(ctx).Info().
		Int64("id", req.Id).
		Msg("users-service: GetUserByID called")

	u, err := s.uc.GetByID(ctx, req.Id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			log.Ctx(ctx).Warn().
				Int64("id", req.Id).
				Msg("users-service: user not found")
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		log.Ctx(ctx).Error().
			Err(err).
			Int64("id", req.Id).
			Msg("users-service: GetByID failed")
		return nil, status.Errorf(codes.Internal, "get user: %v", err)
	}

	return &userspb.GetUserResponse{
		User: &userspb.User{
			Id:       u.Id,
			Email:    u.Email,
			Username: u.Username,
		},
	}, nil
}

// UpdateUser
func (s *GRPCServer) UpdateUser(ctx context.Context, req *userspb.UpdateUserRequest) (*userspb.UpdateUserResponse, error) {
	log.Ctx(ctx).Info().
		Int64("id", req.Id).
		Msg("users-service: UpdateUser called")

	user := &User{
		Id:       req.Id,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password, // usecase сам решит, что делать с паролем
	}

	updated, err := s.uc.Update(ctx, user)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			log.Ctx(ctx).Warn().
				Int64("id", req.Id).
				Msg("users-service: user not found on update")
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		log.Ctx(ctx).Error().
			Err(err).
			Int64("id", req.Id).
			Msg("users-service: Update failed")
		return nil, status.Errorf(codes.Internal, "update user: %v", err)
	}

	return &userspb.UpdateUserResponse{
		User: &userspb.User{
			Id:       updated.Id,
			Email:    updated.Email,
			Username: updated.Username,
		},
	}, nil
}

// DeleteUser
func (s *GRPCServer) DeleteUser(ctx context.Context, req *userspb.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Ctx(ctx).Info().
		Int64("id", req.Id).
		Msg("users-service: DeleteUser called")

	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			log.Ctx(ctx).Warn().
				Int64("id", req.Id).
				Msg("users-service: user not found on delete")
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		log.Ctx(ctx).Error().
			Err(err).
			Int64("id", req.Id).
			Msg("users-service: Delete failed")
		return nil, status.Errorf(codes.Internal, "delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}
