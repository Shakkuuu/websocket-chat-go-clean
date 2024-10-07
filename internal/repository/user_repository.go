package repository

import (
	"context"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/postgres"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/user_mock.go -package=mock_$GOPACKAGE

type UserRepo interface {
	GetAll(ctx context.Context) (*domain.Users, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByName(ctx context.Context, name string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User, id string) error
	Delete(ctx context.Context, id string) error
	NameExists(ctx context.Context, name string) (*bool, error)
}

type userRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) UserRepo {
	return &userRepo{pg}
}

func (r *userRepo) GetAll(ctx context.Context) (*domain.Users, error) {
	var users domain.Users
	err := r.Db.WithContext(ctx).Find(&users).Error
	return &users, err
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.Db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepo) GetByName(ctx context.Context, name string) (*domain.User, error) {
	var user domain.User
	err := r.Db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	return &user, err
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	return r.Db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) Update(ctx context.Context, user *domain.User, id string) error {
	return r.Db.WithContext(ctx).Model(&user).Where("id = ?", id).Updates(&user).Error
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	return r.Db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}).Error
}

func (r *userRepo) NameExists(ctx context.Context, name string) (*bool, error) {
	var exists bool
	err := r.Db.Model(&domain.User{}).Select("count(*) > 0").Where("name = ?", name).Find(&exists).Error
	return &exists, err
}
