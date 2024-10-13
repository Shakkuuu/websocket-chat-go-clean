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

func Test_participatingRoomUsecase_GetAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context)
		want    *domain.ParticipatingRooms
		wantErr bool
	}{
		{
			name: "[正常系] ParticipatingRoom全取得",
			args: args{context.Background()},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context) {
				m.EXPECT().GetAll(ctx).Return(&domain.ParticipatingRooms{domain.ParticipatingRoom{
					ID:        1,
					RoomID:    "1234",
					IsMaster:  true,
					UserID:    "abcd1234",
					User:      domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
					CreatedAt: testTime,
					UpdatedAt: testTime,
				}}, nil)
			},
			want: &domain.ParticipatingRooms{domain.ParticipatingRoom{
				ID:        1,
				RoomID:    "1234",
				IsMaster:  true,
				UserID:    "abcd1234",
				User:      domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
				CreatedAt: testTime,
				UpdatedAt: testTime,
			}},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetAll）",
			args: args{context.Background()},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context) {
				m.EXPECT().GetAll(ctx).Return(&domain.ParticipatingRooms{}, errors.New("test error"))
			},
			want:    &domain.ParticipatingRooms{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			got, err := test.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("participatingRoomUsecase.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_participatingRoomUsecase_GetByUserID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID string)
		want    *domain.ParticipatingRooms
		wantErr bool
	}{
		{
			name: "[正常系] UserID指定でのParticipatingRoom全取得",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID string) {
				m.EXPECT().GetByUserID(ctx, userID).Return(&domain.ParticipatingRooms{domain.ParticipatingRoom{
					ID:        1,
					RoomID:    "1234",
					IsMaster:  true,
					UserID:    "abcd1234",
					User:      domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
					CreatedAt: testTime,
					UpdatedAt: testTime,
				}}, nil)
			},
			want: &domain.ParticipatingRooms{domain.ParticipatingRoom{
				ID:        1,
				RoomID:    "1234",
				IsMaster:  true,
				UserID:    "abcd1234",
				User:      domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
				CreatedAt: testTime,
				UpdatedAt: testTime,
			}},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetByUserID）",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID string) {
				m.EXPECT().GetByUserID(ctx, userID).Return(&domain.ParticipatingRooms{}, errors.New("test error"))
			},
			want:    &domain.ParticipatingRooms{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.userID)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			got, err := test.GetByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.GetByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("participatingRoomUsecase.GetByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_participatingRoomUsecase_GetByRoomID(t *testing.T) {
	type args struct {
		ctx    context.Context
		roomID string
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string)
		want    *domain.ParticipatingRooms
		wantErr bool
	}{
		{
			name: "[正常系] RoomID指定でのParticipatingRoom全取得",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string) {
				m.EXPECT().GetByRoomID(ctx, roomID).Return(&domain.ParticipatingRooms{domain.ParticipatingRoom{
					ID:        1,
					RoomID:    "1234",
					IsMaster:  true,
					UserID:    "abcd1234",
					User:      domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
					CreatedAt: testTime,
					UpdatedAt: testTime,
				}}, nil)
			},
			want: &domain.ParticipatingRooms{domain.ParticipatingRoom{
				ID:        1,
				RoomID:    "1234",
				IsMaster:  true,
				UserID:    "abcd1234",
				User:      domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
				CreatedAt: testTime,
				UpdatedAt: testTime,
			}},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetByRoomID）",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string) {
				m.EXPECT().GetByRoomID(ctx, roomID).Return(&domain.ParticipatingRooms{}, errors.New("test error"))
			},
			want:    &domain.ParticipatingRooms{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.roomID)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			got, err := test.GetByRoomID(tt.args.ctx, tt.args.roomID)
			if (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.GetByRoomID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("participatingRoomUsecase.GetByRoomID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_participatingRoomUsecase_GetByUserIDAndRoomID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		roomID string
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID, roomID string)
		want    *domain.ParticipatingRoom
		wantErr bool
	}{
		{
			name: "[正常系] UserID,RoomID指定でのParticipatingRoom全取得",
			args: args{context.Background(), "abcd1234", "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID, roomID string) {
				m.EXPECT().GetByUserIDAndRoomID(ctx, userID, roomID).Return(&domain.ParticipatingRoom{
					ID:        1,
					RoomID:    "1234",
					IsMaster:  true,
					UserID:    "abcd1234",
					User:      domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
					CreatedAt: testTime,
					UpdatedAt: testTime,
				}, nil)
			},
			want: &domain.ParticipatingRoom{
				ID:        1,
				RoomID:    "1234",
				IsMaster:  true,
				UserID:    "abcd1234",
				User:      domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
				CreatedAt: testTime,
				UpdatedAt: testTime,
			},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetByUserIDAndRoomID）",
			args: args{context.Background(), "abcd1234", "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID, roomID string) {
				m.EXPECT().GetByUserIDAndRoomID(ctx, userID, roomID).Return(&domain.ParticipatingRoom{}, errors.New("test error"))
			},
			want:    &domain.ParticipatingRoom{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.userID, tt.args.roomID)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			got, err := test.GetByUserIDAndRoomID(tt.args.ctx, tt.args.userID, tt.args.roomID)
			if (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.GetByUserIDAndRoomID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("participatingRoomUsecase.GetByUserIDAndRoomID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_participatingRoomUsecase_Create(t *testing.T) {
	type args struct {
		ctx               context.Context
		participatingRoom *domain.ParticipatingRoom
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, participatingRoom *domain.ParticipatingRoom)
		wantErr bool
	}{
		{
			name: "[正常系] ParticipatingRoom作成",
			args: args{context.Background(), &domain.ParticipatingRoom{RoomID: "1234", IsMaster: true, UserID: "abcd1234"}},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, participatingRoom *domain.ParticipatingRoom) {
				m.EXPECT().Create(ctx, participatingRoom).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（Create）",
			args: args{context.Background(), &domain.ParticipatingRoom{RoomID: "1234", IsMaster: true, UserID: "abcd1234"}},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, participatingRoom *domain.ParticipatingRoom) {
				m.EXPECT().Create(ctx, participatingRoom).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.participatingRoom)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			if err := test.Create(tt.args.ctx, tt.args.participatingRoom); (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_participatingRoomUsecase_DeleteByUserID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID string)
		wantErr bool
	}{
		{
			name: "[正常系] UserID指定でのParticipatingRoom削除",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID string) {
				m.EXPECT().DeleteByUserID(ctx, userID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（DeleteByUserID）",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID string) {
				m.EXPECT().DeleteByUserID(ctx, userID).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.userID)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			if err := test.DeleteByUserID(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.DeleteByUserID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_participatingRoomUsecase_DeleteByRoomID(t *testing.T) {
	type args struct {
		ctx    context.Context
		roomID string
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string)
		wantErr bool
	}{
		{
			name: "[正常系] RoomID指定でのParticipatingRoom削除",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string) {
				m.EXPECT().DeleteByRoomID(ctx, roomID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（DeleteByRoomID）",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string) {
				m.EXPECT().DeleteByRoomID(ctx, roomID).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.roomID)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			if err := test.DeleteByRoomID(tt.args.ctx, tt.args.roomID); (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.DeleteByRoomID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_participatingRoomUsecase_DeleteByUserIDAndRoomID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		roomID string
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID, roomID string)
		wantErr bool
	}{
		{
			name: "[正常系] UserID,RoomID指定でのParticipatingRoom削除",
			args: args{context.Background(), "abcd1234", "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID, roomID string) {
				m.EXPECT().DeleteByUserIDAndRoomID(ctx, userID, roomID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（DeleteByUserIDAndRoomID）",
			args: args{context.Background(), "abcd1234", "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, userID, roomID string) {
				m.EXPECT().DeleteByUserIDAndRoomID(ctx, userID, roomID).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.userID, tt.args.roomID)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			if err := test.DeleteByUserIDAndRoomID(tt.args.ctx, tt.args.userID, tt.args.roomID); (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.DeleteByUserIDAndRoomID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_participatingRoomUsecase_GetUsersByRoomID(t *testing.T) {
	type args struct {
		ctx    context.Context
		roomID string
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string)
		want    *domain.Users
		wantErr bool
	}{
		{
			name: "[正常系] RoomID指定でのRoom参加ユーザ一覧取得",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string) {
				m.EXPECT().GetUsersByRoomID(ctx, roomID).Return(&domain.Users{domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime}}, nil)
			},
			want:    &domain.Users{domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime}},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetUsersByRoomID）",
			args: args{context.Background(), "1234"},
			mockFn: func(m *mock_repository.MockParticipatingRoomRepo, ctx context.Context, roomID string) {
				m.EXPECT().GetUsersByRoomID(ctx, roomID).Return(&domain.Users{}, errors.New("test error"))
			},
			want:    &domain.Users{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockParticipatingRoomRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.roomID)

			test := &participatingRoomUsecase{
				repo: mock,
			}
			got, err := test.GetUsersByRoomID(tt.args.ctx, tt.args.roomID)
			if (err != nil) != tt.wantErr {
				t.Errorf("participatingRoomUsecase.GetUsersByRoomID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("participatingRoomUsecase.GetUsersByRoomID() = %v, want %v", got, tt.want)
			}
		})
	}
}
