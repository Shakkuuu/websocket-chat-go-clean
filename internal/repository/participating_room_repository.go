package repository

import (
	"context"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/postgres"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/participating_room_mock.go -package=mock_$GOPACKAGE

type ParticipatingRoomRepo interface {
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

type participatingRoomRepo struct {
	*postgres.Postgres
}

func NewParticipatingRoomRepo(pg *postgres.Postgres) ParticipatingRoomRepo {
	return &participatingRoomRepo{pg}
}

func (r *participatingRoomRepo) GetAll(ctx context.Context) (*domain.ParticipatingRooms, error) {
	var participatingRooms domain.ParticipatingRooms
	err := r.Db.WithContext(ctx).Find(&participatingRooms).Error
	return &participatingRooms, err
}

func (r *participatingRoomRepo) GetByUserID(ctx context.Context, userID string) (*domain.ParticipatingRooms, error) {
	var participatingRooms domain.ParticipatingRooms
	err := r.Db.WithContext(ctx).Where("user_id = ?", userID).Find(&participatingRooms).Error
	return &participatingRooms, err
}

func (r *participatingRoomRepo) GetByRoomID(ctx context.Context, roomID string) (*domain.ParticipatingRooms, error) {
	var participatingRooms domain.ParticipatingRooms
	err := r.Db.WithContext(ctx).Where("room_id = ?", roomID).Find(&participatingRooms).Error
	return &participatingRooms, err
}

func (r *participatingRoomRepo) GetByUserIDAndRoomID(ctx context.Context, userID, roomID string) (*domain.ParticipatingRoom, error) {
	var participatingRoom domain.ParticipatingRoom
	err := r.Db.WithContext(ctx).Where("user_id = ?", userID).Where("room_id = ?", roomID).First(&participatingRoom).Error
	return &participatingRoom, err
}

func (r *participatingRoomRepo) Create(ctx context.Context, participatingRoom *domain.ParticipatingRoom) error {
	return r.Db.WithContext(ctx).Create(participatingRoom).Error
}

func (r *participatingRoomRepo) DeleteByUserID(ctx context.Context, userID string) error {
	return r.Db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.ParticipatingRoom{}).Error
}

func (r *participatingRoomRepo) DeleteByRoomID(ctx context.Context, roomID string) error {
	return r.Db.WithContext(ctx).Where("room_id = ?", roomID).Delete(&domain.ParticipatingRoom{}).Error
}

func (r *participatingRoomRepo) DeleteByUserIDAndRoomID(ctx context.Context, userID, roomID string) error {
	return r.Db.WithContext(ctx).Where("user_id = ?", userID).Where("room_id = ?", roomID).Delete(&domain.ParticipatingRoom{}).Error
}

func (r *participatingRoomRepo) GetUsersByRoomID(ctx context.Context, roomID string) (*domain.Users, error) {
	var participatingRooms domain.ParticipatingRooms
	err := r.Db.WithContext(ctx).Where("room_id = ?", roomID).Preload("User").Find(&participatingRooms).Error

	var users domain.Users
	for _, v := range participatingRooms {
		users = append(users, v.User)
	}
	return &users, err
}
