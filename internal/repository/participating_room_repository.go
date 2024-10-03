package repository

import (
	"context"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/postgres"
)

type ParticipatingRoomRepo struct {
	*postgres.Postgres
}

func NewParticipatingRoomRepo(pg *postgres.Postgres) *ParticipatingRoomRepo {
	return &ParticipatingRoomRepo{pg}
}

func (r *ParticipatingRoomRepo) GetAll(ctx context.Context) (*domain.ParticipatingRooms, error) {
	var participatingRooms domain.ParticipatingRooms
	err := r.Db.WithContext(ctx).Find(&participatingRooms).Error
	return &participatingRooms, err
}

func (r *ParticipatingRoomRepo) GetByUserID(ctx context.Context, userID string) (*domain.ParticipatingRooms, error) {
	var participatingRooms domain.ParticipatingRooms
	err := r.Db.WithContext(ctx).Where("user_id = ?", userID).Find(&participatingRooms).Error
	return &participatingRooms, err
}

func (r *ParticipatingRoomRepo) GetByRoomID(ctx context.Context, roomID string) (*domain.ParticipatingRooms, error) {
	var participatingRooms domain.ParticipatingRooms
	err := r.Db.WithContext(ctx).Where("room_id = ?", roomID).Find(&participatingRooms).Error
	return &participatingRooms, err
}

func (r *ParticipatingRoomRepo) GetByUserIDAndRoomID(ctx context.Context, userID, roomID string) (*domain.ParticipatingRoom, error) {
	var participatingRoom domain.ParticipatingRoom
	err := r.Db.WithContext(ctx).Where("user_id = ?", userID).Where("room_id = ?", roomID).First(&participatingRoom).Error
	return &participatingRoom, err
}

func (r *ParticipatingRoomRepo) Create(ctx context.Context, participatingRoom *domain.ParticipatingRoom) error {
	return r.Db.WithContext(ctx).Create(participatingRoom).Error
}

func (r *ParticipatingRoomRepo) DeleteByUserID(ctx context.Context, userID string) error {
	return r.Db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.ParticipatingRoom{}).Error
}

func (r *ParticipatingRoomRepo) DeleteByRoomID(ctx context.Context, roomID string) error {
	return r.Db.WithContext(ctx).Where("room_id = ?", roomID).Delete(&domain.ParticipatingRoom{}).Error
}

func (r *ParticipatingRoomRepo) DeleteByUserIDAndRoomID(ctx context.Context, userID, roomID string) error {
	return r.Db.WithContext(ctx).Where("user_id = ?", userID).Where("room_id = ?", roomID).Delete(&domain.ParticipatingRoom{}).Error
}

func (r *ParticipatingRoomRepo) GetUsersByRoomID(ctx context.Context, roomID string) (*domain.Users, error) {
	var participatingRooms domain.ParticipatingRooms
	err := r.Db.WithContext(ctx).Where("room_id = ?", roomID).Preload("User").Find(&participatingRooms).Error

	var users domain.Users
	for _, v := range participatingRooms {
		users = append(users, v.User)
	}
	return &users, err
}
