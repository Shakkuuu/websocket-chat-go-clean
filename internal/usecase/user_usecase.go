package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/repository"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/ulid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/user_mock.go -package=mock_$GOPACKAGE

type UserUsecase interface {
	GetAll(ctx context.Context) (*domain.Users, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByName(ctx context.Context, name string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User, id string) error
	Delete(ctx context.Context, id string) error
	NameExists(ctx context.Context, name string) (*bool, error)
}

type userUsecase struct {
	repo repository.UserRepo
}

func NewUserUsecase(repo repository.UserRepo) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetAll(ctx context.Context) (*domain.Users, error) {
	return u.repo.GetAll(ctx)
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUsecase) GetByName(ctx context.Context, name string) (*domain.User, error) {
	return u.repo.GetByName(ctx, name)
}

func (u *userUsecase) Create(ctx context.Context, user *domain.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	exists, err := u.repo.NameExists(ctx, user.Name)
	if err != nil {
		return err
	}
	if *exists {
		return errors.New("その名前は既に登録されています。")
	}

	user.ID = ulid.NewULID()

	hp, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hp)

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return u.repo.Create(ctx, user)
}

func (u *userUsecase) Update(ctx context.Context, user *domain.User, id string) error {
	now := time.Now()
	user.UpdatedAt = now
	return u.repo.Update(ctx, user, id)
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *userUsecase) NameExists(ctx context.Context, name string) (*bool, error) {
	return u.repo.NameExists(ctx, name)
}
