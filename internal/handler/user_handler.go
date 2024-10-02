package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/usecase"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/session"
)

type UserHandler struct {
	userUsecase              usecase.UserUsecase
	participatingRoomUsecase usecase.ParticipatingRoomUsecase
	roomUsecase              usecase.RoomUsecase
	templates                *template.Template
	session                  *session.Sessions
}

func NewUserHandler(
	usecase usecase.UserUsecase,
	participatingRoomUsecase usecase.ParticipatingRoomUsecase,
	roomUsecase usecase.RoomUsecase,
	s *session.Sessions,
) *UserHandler {
	templates := template.Must(template.ParseGlob("internal/handler/templates/*.html"))
	return &UserHandler{
		userUsecase:              usecase,
		participatingRoomUsecase: participatingRoomUsecase,
		roomUsecase:              roomUsecase,
		templates:                templates,
		session:                  s,
	}
}

// ユーザー情報のページ
func (h *UserHandler) Menu(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// セッション読み取り
		userID, _, err := h.session.GetUserData(r)
		if err != nil {
			log.Printf("session.GetUserData error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "再ログインしてください"

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		user, err := h.userUsecase.GetByID(ctx, userID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// メッセージをテンプレートに渡す
		var data Data
		data.Message = user.Name + "さん、こんにちは。"

		data.Name = user.Name

		err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// ユーザー削除
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// セッション読み取り
		userID, _, err := h.session.GetUserData(r)
		if err != nil {
			log.Printf("session.GetUserData error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "再ログインしてください"

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// セッションのユーザー取得
		user, err := h.userUsecase.GetByID(ctx, userID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// セッション削除
		h.session.Delete(r, w)

		// ユーザーが作成したRoomの削除
		prooms, err := h.participatingRoomUsecase.GetByUserID(ctx, user.ID)
		if err != nil {
			log.Printf("participatingRoomUsecase.GetByUserID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		for _, proom := range *prooms {
			if !proom.IsMaster {
				continue
			}

			// ユーザーの参加中ルームリストからも削除
			err := h.participatingRoomUsecase.DeleteByRoomID(ctx, proom.RoomID)
			if err != nil {
				log.Printf("participatingRoomUsecase.DeleteByRoomID error: %v\n", err)
				// メッセージをテンプレートに渡す
				var data Data
				data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

				err := h.templates.ExecuteTemplate(w, "usermenu.html", data)
				if err != nil {
					log.Printf("templates.ExecuteTemplate error:%v\n", err)
					http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
					return
				}
				return
			}

			// 部屋削除
			err = h.roomUsecase.Delete(ctx, proom.RoomID)
			if err != nil {
				log.Printf("roomUsecase.Delete error: %v\n", err)
				// メッセージをテンプレートに渡す
				var data Data
				data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

				err := h.templates.ExecuteTemplate(w, "usermenu.html", data)
				if err != nil {
					log.Printf("templates.ExecuteTemplate error:%v\n", err)
					http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
					return
				}
				return
			}
			deleteRoom(proom.RoomID)
		}

		// ユーザーの参加中ルームリストを削除
		err = h.participatingRoomUsecase.DeleteByUserID(ctx, user.ID)
		if err != nil {
			log.Printf("participatingRoomUsecase.DeleteByUserID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// ユーザー削除
		err = h.userUsecase.Delete(ctx, userID)
		if err != nil {
			log.Printf("userUsecase.Delete error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// メッセージをテンプレートに渡す
		var data Data
		data.Message = "ユーザーを削除しました。"
		err = h.templates.ExecuteTemplate(w, "signup.html", data)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
		return
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// パスワード変更
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		r.ParseForm()
		oldpassword := r.FormValue("oldpassword")
		password := r.FormValue("password")
		checkpass := r.FormValue("checkpassword")

		// セッション読み取り
		userID, _, err := h.session.GetUserData(r)
		if err != nil {
			log.Printf("session.GetUserData error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "再ログインしてください"

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// セッションのユーザー取得
		user, err := h.userUsecase.GetByID(ctx, userID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByID error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		if oldpassword == "" || password == "" || checkpass == "" {
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "入力されていない項目があります。"

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		if password != checkpass {
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "確認用再入力パスワードが一致していません。"

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// ハッシュ化されたパスワードの解読と一致確認
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldpassword))
		if err != nil {
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "現在のパスワードが違います"

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// hashpass, err := model.HashPass(password)
		hp, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("bcrypt.GenerateFromPassword error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "パスワードのハッシュに失敗しました。"

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		hashpass := string(hp)

		updatedUser := domain.User{
			ID:        user.ID,
			Name:      user.Name,
			Password:  hashpass,
			CreatedAt: user.CreatedAt,
		}

		// ユーザー更新
		err = h.userUsecase.Update(ctx, &updatedUser, userID)
		if err != nil {
			log.Printf("userUsecase.Update error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err = h.templates.ExecuteTemplate(w, "usermenu.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// 再ログイン用に一度セッション削除
		h.session.Delete(r, w)

		// メッセージをテンプレートに渡す
		var data Data
		data.Message = "パスワードを更新しました。再ログインしてください。"

		err = h.templates.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
		return
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// Signup処理
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.templates.ExecuteTemplate(w, "signup.html", nil)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// POSTされたものをFormから受け取り
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		checkpass := r.FormValue("checkpassword")

		if username == "" || password == "" || checkpass == "" {
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "入力されていない項目があります。"

			err := h.templates.ExecuteTemplate(w, "signup.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		if password != checkpass {
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "確認用再入力パスワードが一致していません。"

			err := h.templates.ExecuteTemplate(w, "signup.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		user := domain.User{
			Name:     username,
			Password: password,
		}

		// ユーザー追加
		err := h.userUsecase.Create(ctx, &user)
		if err != nil {
			log.Printf("model.AddUser error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err := h.templates.ExecuteTemplate(w, "signup.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// メッセージをテンプレートに渡す
		var data Data
		data.Message = "登録が完了しました。ログインしてください。"

		err = h.templates.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// Login処理
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// POSTされたものをFormから受け取り
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "入力されていない項目があります。"

			err := h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// 登録されているユーザー取得
		user, err := h.userUsecase.GetByName(ctx, username)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByName error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "ユーザーが見つかりませんでした。"

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByName error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = fmt.Sprintf("データベースとの接続に失敗しました。(%v)", err)

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// ハッシュ化されたパスワードの解読と一致確認
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "パスワードが違います"

			err = h.templates.ExecuteTemplate(w, "login.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// ログイン成功時の処理

		// メッセージをテンプレートに渡す
		var data Data
		data.Message = "ログインに成功しました。"

		h.session.Set(r, w, user.ID, user.Name)

		err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// Logout処理
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// セッション削除
		h.session.Delete(r, w)

		// メッセージをテンプレートに渡す
		var data Data
		data.Message = "ログアウトしました。ログインしてください。"

		err := h.templates.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// 自身のユーザー名を返す
func (h *UserHandler) GetUserName(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var sentuser SentUser

		// セッション読み取り
		_, un, err := h.session.GetUserData(r)
		if err != nil {
			log.Printf("SessionToGetName error: %v\n", err)
			log.Println("セッションが見つかりませんでした")
			http.Error(w, "セッションが見つかりませんでした", http.StatusNotFound)
			return
		}

		sentuser.Name = un

		// jsonに変換
		sentjson, err := json.Marshal(sentuser)
		if err != nil {
			log.Printf("json.Marshal error: %v\n", err)
			http.Error(w, "json.Marshal error", http.StatusInternalServerError)
			return
		}

		// jsonで送信
		w.Header().Set("Content-Type", "application/json")
		w.Write(sentjson)
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}
