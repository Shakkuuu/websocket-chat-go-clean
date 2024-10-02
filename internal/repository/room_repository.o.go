package repository

import (
	"context"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/postgres"
)

type RoomRepo struct {
	*postgres.Postgres
}

func NewRoomRepo(pg *postgres.Postgres) *RoomRepo {
	return &RoomRepo{pg}
}

func (r *RoomRepo) GetAll(ctx context.Context) (*domain.Rooms, error) {
	var rooms domain.Rooms
	err := r.Db.WithContext(ctx).Find(&rooms).Error
	return &rooms, err
}

func (r *RoomRepo) Create(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	err := r.Db.WithContext(ctx).Create(room).Error
	return room, err
}

func (r *RoomRepo) Delete(ctx context.Context, id string) error {
	return r.Db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Room{}).Error
}

func (r *RoomRepo) IDExists(ctx context.Context, id string) (*bool, error) {
	var exists bool
	err := r.Db.Model(&domain.Room{}).Select("count(*) > 0").Where("id = ?", id).Find(&exists).Error
	return &exists, err
}
