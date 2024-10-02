package repository

import (
	"context"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) GetAll(ctx context.Context) (*domain.Users, error) {
	var users domain.Users
	err := r.Db.WithContext(ctx).Find(&users).Error
	return &users, err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.Db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepo) GetByName(ctx context.Context, name string) (*domain.User, error) {
	var user domain.User
	err := r.Db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	return &user, err
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	return r.Db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User, id string) error {
	return r.Db.WithContext(ctx).Model(&user).Where("id = ?", id).Updates(&user).Error
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	return r.Db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}).Error
}

func (r *UserRepo) NameExists(ctx context.Context, name string) (*bool, error) {
	var exists bool
	err := r.Db.Model(&domain.User{}).Select("count(*) > 0").Where("name = ?", name).Find(&exists).Error
	return &exists, err
}
