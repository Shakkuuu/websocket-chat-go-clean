package usecase

import (
	"context"
	"time"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/repository"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/participating_room_mock.go -package=mock_$GOPACKAGE

type ParticipatingRoomUsecase interface {
	GetAll(ctx context.Context) (*domain.ParticipatingRooms, error)
	GetByUserID(ctx context.Context, userID string) (*domain.ParticipatingRooms, error)
	GetByRoomID(ctx context.Context, roomID string) (*domain.ParticipatingRooms, error)
	GetByUserIDAndRoomID(ctx context.Context, userID, roomID string) (*domain.ParticipatingRoom, error)
	Create(ctx context.Context, participatingRoom *domain.ParticipatingRoom) error
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteByRoomID(ctx context.Context, roomID string) error
	DeleteByUserIDAndRoomID(ctx context.Context, userID, roomID string) error
	GetUsersByRoomID(ctx context.Context, roomID string) (*domain.Users, error)
}

type participatingRoomUsecase struct {
	repo repository.ParticipatingRoomRepo
}

func NewParticipatingRoomUsecase(repo repository.ParticipatingRoomRepo) ParticipatingRoomUsecase {
	return &participatingRoomUsecase{repo: repo}
}

func (u *participatingRoomUsecase) GetAll(ctx context.Context) (*domain.ParticipatingRooms, error) {
	return u.repo.GetAll(ctx)
}

func (u *participatingRoomUsecase) GetByUserID(ctx context.Context, userID string) (*domain.ParticipatingRooms, error) {
	return u.repo.GetByUserID(ctx, userID)
}

func (u *participatingRoomUsecase) GetByRoomID(ctx context.Context, roomID string) (*domain.ParticipatingRooms, error) {
	return u.repo.GetByRoomID(ctx, roomID)
}

func (u *participatingRoomUsecase) GetByUserIDAndRoomID(ctx context.Context, userID, roomID string) (*domain.ParticipatingRoom, error) {
	return u.repo.GetByUserIDAndRoomID(ctx, userID, roomID)
}

func (u *participatingRoomUsecase) Create(ctx context.Context, participatingRoom *domain.ParticipatingRoom) error {
	now := time.Now()
	participatingRoom.CreatedAt = now
	participatingRoom.UpdatedAt = now

	return u.repo.Create(ctx, participatingRoom)
}

func (u *participatingRoomUsecase) DeleteByUserID(ctx context.Context, userID string) error {
	return u.repo.DeleteByUserID(ctx, userID)
}

func (u *participatingRoomUsecase) DeleteByRoomID(ctx context.Context, roomID string) error {
	return u.repo.DeleteByRoomID(ctx, roomID)
}

func (u *participatingRoomUsecase) DeleteByUserIDAndRoomID(ctx context.Context, userID, roomID string) error {
	return u.repo.DeleteByUserIDAndRoomID(ctx, userID, roomID)
}

func (u *participatingRoomUsecase) GetUsersByRoomID(ctx context.Context, roomID string) (*domain.Users, error) {
	return u.repo.GetUsersByRoomID(ctx, roomID)
}
