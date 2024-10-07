package repository

import (
	"context"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/postgres"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/room_mock.go -package=mock_$GOPACKAGE

type RoomRepo interface {
	GetAll(ctx context.Context) (*domain.Rooms, error)
	Create(ctx context.Context, room *domain.Room) (*domain.Room, error)
	Delete(ctx context.Context, id string) error
	IDExists(ctx context.Context, id string) (*bool, error)
}

type roomRepo struct {
	*postgres.Postgres
}

func NewRoomRepo(pg *postgres.Postgres) RoomRepo {
	return &roomRepo{pg}
}

func (r *roomRepo) GetAll(ctx context.Context) (*domain.Rooms, error) {
	var rooms domain.Rooms
	err := r.Db.WithContext(ctx).Find(&rooms).Error
	return &rooms, err
}

func (r *roomRepo) Create(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	err := r.Db.WithContext(ctx).Create(room).Error
	return room, err
}

func (r *roomRepo) Delete(ctx context.Context, id string) error {
	return r.Db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Room{}).Error
}

func (r *roomRepo) IDExists(ctx context.Context, id string) (*bool, error) {
	var exists bool
	err := r.Db.Model(&domain.Room{}).Select("count(*) > 0").Where("id = ?", id).Find(&exists).Error
	return &exists, err
}
