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

func Test_userUsecase_GetAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockUserRepo, ctx context.Context)
		want    *domain.Users
		wantErr bool
	}{
		{
			name: "[正常系] ユーザー全取得",
			args: args{context.Background()},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context) {
				m.EXPECT().GetAll(ctx).Return(&domain.Users{domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime}}, nil)
			},
			want:    &domain.Users{domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime}},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetAll）",
			args: args{context.Background()},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context) {
				m.EXPECT().GetAll(ctx).Return(&domain.Users{}, errors.New("test error"))
			},
			want:    &domain.Users{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockUserRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx)

			test := &userUsecase{
				repo: mock,
			}
			got, err := test.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userUsecase.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userUsecase_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockUserRepo, ctx context.Context, id string)
		want    *domain.User
		wantErr bool
	}{
		{
			name: "[正常系] ID指定でのユーザー取得",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, id string) {
				m.EXPECT().GetByID(ctx, id).Return(&domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime}, nil)
			},
			want:    &domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetByID）",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, id string) {
				m.EXPECT().GetByID(ctx, id).Return(&domain.User{}, errors.New("test error"))
			},
			want:    &domain.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockUserRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.id)

			test := &userUsecase{
				repo: mock,
			}
			got, err := test.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userUsecase.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userUsecase_GetByName(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockUserRepo, ctx context.Context, name string)
		want    *domain.User
		wantErr bool
	}{
		{
			name: "[正常系] Name指定でのユーザー取得",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				m.EXPECT().GetByName(ctx, name).Return(&domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime}, nil)
			},
			want:    &domain.User{ID: "abcd1234", Name: "testName", Password: "testPass", CreatedAt: testTime, UpdatedAt: testTime},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（GetByName）",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				m.EXPECT().GetByName(ctx, name).Return(&domain.User{}, errors.New("test error"))
			},
			want:    &domain.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockUserRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.name)

			test := &userUsecase{
				repo: mock,
			}
			got, err := test.GetByName(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userUsecase.GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userUsecase_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *domain.User
	}
	tests := []struct {
		name    string
		args    args
		mockFn1 func(m *mock_repository.MockUserRepo, ctx context.Context, name string)
		mockFn2 func(m *mock_repository.MockUserRepo, ctx context.Context, user *domain.User)
		wantErr bool
	}{
		{
			name: "[正常系] ユーザー作成",
			args: args{context.Background(), &domain.User{Name: "testName", Password: "p@ssw0rd"}},
			mockFn1: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				exists := false
				m.EXPECT().NameExists(ctx, name).Return(&exists, nil)
			},
			mockFn2: func(m *mock_repository.MockUserRepo, ctx context.Context, user *domain.User) {
				m.EXPECT().Create(ctx, user).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "[異常系] バリデーション失敗（Nameが空）",
			args:    args{context.Background(), &domain.User{Name: "", Password: "p@ssw0rd"}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name:    "[異常系] バリデーション失敗（Passwordが空）",
			args:    args{context.Background(), &domain.User{Name: "testName", Password: ""}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name: "[異常系] バリデーション失敗（Nameが100文字より大きい）",
			args: args{context.Background(), &domain.User{
				Name:     "0123456789112345678921234567893123456789412345678951234567896123456789712345678981234567899123456789n",
				Password: "p@ssw0rd",
			}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name:    "[異常系] バリデーション失敗（Passwordが8文字未満）",
			args:    args{context.Background(), &domain.User{Name: "testName", Password: "1nval1d"}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name: "[異常系] バリデーション失敗（Passwordが100文字より大きい）",
			args: args{context.Background(), &domain.User{
				Name:     "testName",
				Password: "0123456789112345678921234567893123456789412345678951234567896123456789712345678981234567899123456789p",
			}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name:    "[異常系] バリデーション失敗（Passwordが半角英字のみ）",
			args:    args{context.Background(), &domain.User{Name: "testName", Password: "invalidPass"}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name:    "[異常系] バリデーション失敗（Passwordが半角数字のみ）",
			args:    args{context.Background(), &domain.User{Name: "testName", Password: "01234567"}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name:    "[異常系] バリデーション失敗（Passwordが記号のみ）",
			args:    args{context.Background(), &domain.User{Name: "testName", Password: "!@#$%^&*"}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name:    "[異常系] バリデーション失敗（Passwordに許可されていない記号が使用されている）",
			args:    args{context.Background(), &domain.User{Name: "testName", Password: "p@ss,w0rd"}},
			mockFn1: nil,
			mockFn2: nil,
			wantErr: true,
		},
		{
			name: "[異常系] 名前がすでに存在している",
			args: args{context.Background(), &domain.User{Name: "existsName", Password: "p@ssw0rd"}},
			mockFn1: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				exists := true
				m.EXPECT().NameExists(ctx, name).Return(&exists, nil)
			},
			mockFn2: nil,
			wantErr: true,
		},
		{
			name: "[異常系] DB処理失敗（NameExists）",
			args: args{context.Background(), &domain.User{Name: "testName", Password: "p@ssw0rd"}},
			mockFn1: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				exists := false
				m.EXPECT().NameExists(ctx, name).Return(&exists, errors.New("test error"))
			},
			mockFn2: nil,
			wantErr: true,
		},
		{
			name: "[異常系] DB処理失敗（Create）",
			args: args{context.Background(), &domain.User{Name: "testName", Password: "p@ssw0rd"}},
			mockFn1: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				exists := false
				m.EXPECT().NameExists(ctx, name).Return(&exists, nil)
			},
			mockFn2: func(m *mock_repository.MockUserRepo, ctx context.Context, user *domain.User) {
				m.EXPECT().Create(ctx, user).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockUserRepo(ctrl)

			if tt.mockFn1 != nil {
				tt.mockFn1(mock, tt.args.ctx, tt.args.user.Name)
			}
			if tt.mockFn2 != nil {
				tt.mockFn2(mock, tt.args.ctx, tt.args.user)
			}

			test := &userUsecase{
				repo: mock,
			}
			if err := test.Create(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userUsecase_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *domain.User
		id   string
	}
	testTime := time.Now()
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockUserRepo, ctx context.Context, user *domain.User, id string)
		wantErr bool
	}{
		{
			name: "[正常系] ユーザー更新",
			args: args{context.Background(), &domain.User{ID: "abcd1234", Name: "testName", Password: "p@ssw0rd", CreatedAt: testTime, UpdatedAt: testTime}, "abcd1234"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, user *domain.User, id string) {
				m.EXPECT().Update(ctx, user, id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（Update）",
			args: args{context.Background(), &domain.User{ID: "abcd1234", Name: "testName", Password: "p@ssw0rd", CreatedAt: testTime, UpdatedAt: testTime}, "abcd1234"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, user *domain.User, id string) {
				m.EXPECT().Update(ctx, user, id).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockUserRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.user, tt.args.id)

			test := &userUsecase{
				repo: mock,
			}
			if err := test.Update(tt.args.ctx, tt.args.user, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userUsecase_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockUserRepo, ctx context.Context, id string)
		wantErr bool
	}{
		{
			name: "[正常系] ユーザー削除",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, id string) {
				m.EXPECT().Delete(ctx, id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（Delete）",
			args: args{context.Background(), "abcd1234"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, id string) {
				m.EXPECT().Delete(ctx, id).Return(errors.New("test error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockUserRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.id)

			test := &userUsecase{
				repo: mock,
			}
			if err := test.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userUsecase_NameExists(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	wantTrue := true
	wantFalse := false
	tests := []struct {
		name    string
		args    args
		mockFn  func(m *mock_repository.MockUserRepo, ctx context.Context, name string)
		want    *bool
		wantErr bool
	}{
		{
			name: "[正常系] 名前存在確認（あった）",
			args: args{context.Background(), "testName"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				exists := true
				m.EXPECT().NameExists(ctx, name).Return(&exists, nil)
			},
			want:    &wantTrue,
			wantErr: false,
		},
		{
			name: "[正常系] 名前存在確認（なかった）",
			args: args{context.Background(), "testName"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				exists := false
				m.EXPECT().NameExists(ctx, name).Return(&exists, nil)
			},
			want:    &wantFalse,
			wantErr: false,
		},
		{
			name: "[異常系] DB処理失敗（NameExists）",
			args: args{context.Background(), "testName"},
			mockFn: func(m *mock_repository.MockUserRepo, ctx context.Context, name string) {
				var exists bool
				m.EXPECT().NameExists(ctx, name).Return(&exists, errors.New("test error"))
			},
			want:    &wantFalse,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_repository.NewMockUserRepo(ctrl)

			tt.mockFn(mock, tt.args.ctx, tt.args.name)

			test := &userUsecase{
				repo: mock,
			}
			got, err := test.NameExists(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("userUsecase.NameExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != *tt.want {
				t.Errorf("userUsecase.NameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
