package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	mock_repository "github.com/Shakkuuu/websocket-chat-go-clean/internal/mock/repository"
	"go.uber.org/mock/gomock"
)

func Test_roomUsecase_GetAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockRoomRepo, ctx context.Context)
		want    *domain.Rooms
		wantErr bool
	}{
		{
			name: "[正常系] Room全取得",
			args: args{context.Background()},
			mockFn: func(m *mock_repository.MockRoomRepo, ctx context.Context) {
				m.EXPECT().GetAll(ctx).Return(&domain.Rooms{domain.Room{ID: "1234", CreatedAt: testTime, UpdatedAt: testTime}}, nil)
			},
			want:    &domain.Rooms{domain.Room{ID: "1234", CreatedAt: testTime, UpdatedAt: testTime}},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetAll）",
			args: args{context.Background()},
			mockFn: func(m *mock_repository.MockRoomRepo, ctx context.Context) {
				m.EXPECT().GetAll(ctx).Return(&domain.Rooms{}, errors.New("test error"))
			},
			want:    &domain.Rooms{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx)

			test := &roomUsecase{
				repo: mock,
			}
			got, err := test.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("roomUsecase.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roomUsecase.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roomUsecase_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		room *domain.Room
	}
	testTime := time.Now()
	tests := []struct {
		name        string
		args        args
		mockFn1     func(m *mock_repository.MockRoomRepo, ctx context.Context, id string)
		againMockFn func(m *mock_repository.MockRoomRepo, ctx context.Context, id string)
		mockFn2     func(m *mock_repository.MockRoomRepo, ctx context.Context, room *domain.Room)
		want        *domain.Room
		wantErr     bool
	}{
		{
			name: "[正常系] Room作成",
			args: args{context.Background(), &domain.Room{}},
			mockFn1: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				exists := false
				m.EXPECT().IDExists(ctx, id).Return(&exists, nil)
			},
			againMockFn: nil,
			mockFn2: func(m *mock_repository.MockRoomRepo, ctx context.Context, room *domain.Room) {
				m.EXPECT().Create(ctx, room).Return(&domain.Room{ID: "1234", CreatedAt: testTime, UpdatedAt: testTime}, nil)
			},
			want:    &domain.Room{ID: "1234", CreatedAt: testTime, UpdatedAt: testTime},
			wantErr: false,
		},
		{
			name: "[正常系] すでにIDが存在しており再度ID生成処理が走る場合",
			args: args{context.Background(), &domain.Room{}},
			mockFn1: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				exists := true
				m.EXPECT().IDExists(ctx, id).Return(&exists, nil)
			},
			againMockFn: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				exists := false
				m.EXPECT().IDExists(ctx, id).Return(&exists, nil)
			},
			mockFn2: func(m *mock_repository.MockRoomRepo, ctx context.Context, room *domain.Room) {
				m.EXPECT().Create(ctx, room).Return(&domain.Room{ID: "5678", CreatedAt: testTime, UpdatedAt: testTime}, nil)
			},
			want:    &domain.Room{ID: "5678", CreatedAt: testTime, UpdatedAt: testTime},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（IDExists）",
			args: args{context.Background(), &domain.Room{}},
			mockFn1: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				var exists bool
				m.EXPECT().IDExists(ctx, id).Return(&exists, errors.New("test error"))
			},
			againMockFn: nil,
			mockFn2:     nil,
			want:        nil,
			wantErr:     true,
		},
		{
			name: "[異常系] DB処理失敗（Create）",
			args: args{context.Background(), &domain.Room{}},
			mockFn1: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				exists := false
				m.EXPECT().IDExists(ctx, id).Return(&exists, nil)
			},
			againMockFn: nil,
			mockFn2: func(m *mock_repository.MockRoomRepo, ctx context.Context, room *domain.Room) {
				m.EXPECT().Create(ctx, room).Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockRoomRepo(ctrl)

			if tt.mockFn1 != nil {
				tt.mockFn1(mock, tt.args.ctx, tt.args.room.ID)
			}
			if tt.againMockFn != nil {
				tt.againMockFn(mock, tt.args.ctx, tt.args.room.ID)
			}
			if tt.mockFn2 != nil {
				tt.mockFn2(mock, tt.args.ctx, tt.args.room)
			}

			test := &roomUsecase{
				repo: mock,
			}
			got, err := test.Create(tt.args.ctx, tt.args.room)
			if (err != nil) != tt.wantErr {
				t.Errorf("roomUsecase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roomUsecase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roomUsecase_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockRoomRepo, ctx context.Context, id string)
		wantErr bool
	}{
		{
			name: "[正常系] Room削除",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				m.EXPECT().Delete(ctx, id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（Delete）",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				m.EXPECT().Delete(ctx, id).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.id)

			test := &roomUsecase{
				repo: mock,
			}
			if err := test.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("roomUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_roomUsecase_IDExists(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	wantTrue := true
	wantFalse := false
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockRoomRepo, ctx context.Context, id string)
		want    *bool
		wantErr bool
	}{
		{
			name: "[正常系] RoomID存在確認（あった）",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				exists := true
				m.EXPECT().IDExists(ctx, id).Return(&exists, nil)
			},
			want:    &wantTrue,
			wantErr: false,
		},
		{
			name: "[正常系] RoomID存在確認（なかった）",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				exists := false
				m.EXPECT().IDExists(ctx, id).Return(&exists, nil)
			},
			want:    &wantFalse,
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（IDExists）",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockRoomRepo, ctx context.Context, id string) {
				var exists bool
				m.EXPECT().IDExists(ctx, id).Return(&exists, errors.New("test error"))
			},
			want:    &wantFalse,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.id)

			test := &roomUsecase{
				repo: mock,
			}
			got, err := test.IDExists(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("roomUsecase.IDExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != *tt.want {
				t.Errorf("roomUsecase.IDExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
