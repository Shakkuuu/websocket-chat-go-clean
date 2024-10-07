package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/repository"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/room_mock.go -package=mock_$GOPACKAGE

type RoomUsecase interface {
	GetAll(ctx context.Context) (*domain.Rooms, error)
	Create(ctx context.Context, user *domain.Room) (*domain.Room, error)
	Delete(ctx context.Context, id string) error
	IDExists(ctx context.Context, id string) (*bool, error)
}

type roomUsecase struct {
	repo repository.RoomRepo
}

func NewRoomUsecase(repo repository.RoomRepo) RoomUsecase {
	return &roomUsecase{repo: repo}
}

func (u *roomUsecase) GetAll(ctx context.Context) (*domain.Rooms, error) {
	return u.repo.GetAll(ctx)
}

func (u *roomUsecase) Create(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))

	var roomID string
	var exists *bool
	var err error

	for {
		roomID = fmt.Sprintf("%04d", ran.Intn(10000))

		exists, err = u.repo.IDExists(ctx, room.ID)
		if err != nil {
			return nil, err
		}

		if !*exists {
			break
		}
	}

	room.ID = roomID

	now := time.Now()
	room.CreatedAt = now
	room.UpdatedAt = now

	return u.repo.Create(ctx, room)
}

func (u *roomUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *roomUsecase) IDExists(ctx context.Context, id string) (*bool, error) {
	return u.repo.IDExists(ctx, id)
}
