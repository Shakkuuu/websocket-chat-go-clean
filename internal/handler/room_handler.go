package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/usecase"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/session"
	"gorm.io/gorm"
)

type RoomHandler struct {
	userUsecase              usecase.UserUsecase
	participatingRoomUsecase usecase.ParticipatingRoomUsecase
	roomUsecase              usecase.RoomUsecase
	templates                *template.Template
	session                  *session.Sessions
}

func NewRoomHandler(
	usecase usecase.UserUsecase,
	participatingRoomUsecase usecase.ParticipatingRoomUsecase,
	roomUsecase usecase.RoomUsecase,
	s *session.Sessions,
) *RoomHandler {
	templates := template.Must(template.ParseGlob("internal/handler/templates/*.html"))
	roomInit(roomUsecase)
	return &RoomHandler{
		userUsecase:              usecase,
		participatingRoomUsecase: participatingRoomUsecase,
		roomUsecase:              roomUsecase,
		templates:                templates,
		session:                  s,
	}
}

// roomtopページの表示
func (h *RoomHandler) Top(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.templates.ExecuteTemplate(w, "roomtop.html", nil)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
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

			err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
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

		// Room作成
		room, err := h.roomUsecase.Create(ctx, &domain.Room{})
		if err != nil {
			log.Printf("roomUsecase.Create error: %v\n", err)
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
		createRoom(room.ID)

		// 参加中のルーム一覧にMasterとして追加
		proom := domain.ParticipatingRoom{
			RoomID:   room.ID,
			IsMaster: true,
			UserID:   user.ID,
		}
		err = h.participatingRoomUsecase.Create(ctx, &proom)
		if err != nil {
			log.Printf("participatingRoomUsecase.Create error: %v\n", err)
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
		data.Message = "ルーム " + room.ID + " が作成されました。"

		err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
		if err != nil {
			log.Printf("templates.ExecuteTemplate error:%v\n", err)
			http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
			return
		}
		return
	case http.MethodHead:
		fmt.Fprintln(w, "Thank you monitor.")
	default:
		fmt.Fprintln(w, "Method not allowed")
		http.Error(w, "そのメソッドは許可されていません。", http.StatusMethodNotAllowed)
		return
	}
}

// Room内のページ
func (h *RoomHandler) Room(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// クエリ読み取り
		r.ParseForm()
		roomid := r.URL.Query().Get("roomid")

		intRoomID, err := strconv.Atoi(roomid)
		if err != nil {
			log.Printf("strconv.Atoi error: %v\n", err)
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "ルームIDの形式が正しくありません。"

			err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}
		if intRoomID < 1 || 9999 < intRoomID {
			log.Println("ルームIDの範囲外です。")
			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "ルームIDの範囲外です。"

			err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// Roomが存在するか確認
		exists, err := h.roomUsecase.IDExists(ctx, roomid)
		if err != nil {
			log.Printf("roomUsecase.IDExists error: %v\n", err)
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
		if !*exists { // 指定した部屋が存在していなかったら
			log.Println("This room was not found")

			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "そのIDのルームは見つかりませんでした。"

			err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		err = h.templates.ExecuteTemplate(w, "room.html", nil)
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

// Room削除
func (h *RoomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		// クエリ読み取り
		r.ParseForm()
		roomid := r.URL.Query().Get("roomid")

		// Roomが存在するか確認
		exists, err := h.roomUsecase.IDExists(ctx, roomid)
		if err != nil {
			log.Printf("roomUsecase.IDExists error: %v\n", err)
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
		if !*exists { // 指定した部屋が存在していなかったら
			log.Println("This room was not found")

			// メッセージをテンプレートに渡す
			var data Data
			data.Message = "そのIDのルームは見つかりませんでした。"

			err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

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

			err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
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

		proom, err := h.participatingRoomUsecase.GetByUserIDAndRoomID(ctx, user.ID, roomid)
		if err != nil {
			log.Printf("participatingRoomUsecase.GetByUserIDAndRoomID error: %v\n", err)
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

		// 部屋の作成者ではない場合は、部屋から離脱
		if !proom.IsMaster {
			// 部屋離脱
			err := h.participatingRoomUsecase.DeleteByUserIDAndRoomID(ctx, user.ID, roomid)
			if err != nil {
				log.Printf("participatingRoomUsecase.DeleteByUserIDAndRoomID error: %v\n", err)
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
			data.Message = "部屋を離脱しました。"

			err = h.templates.ExecuteTemplate(w, "roomtop.html", data)
			if err != nil {
				log.Printf("templates.ExecuteTemplate error:%v\n", err)
				http.Error(w, "ページの表示に失敗しました。", http.StatusInternalServerError)
				return
			}
			return
		}

		// ユーザーの参加中ルームリストからも削除
		err = h.participatingRoomUsecase.DeleteByRoomID(ctx, roomid)
		if err != nil {
			log.Printf("participatingRoomUsecase.DeleteByRoomID error: %v\n", err)
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

		// 部屋削除
		err = h.roomUsecase.Delete(ctx, roomid)
		if err != nil {
			log.Printf("roomUsecase.Delete error: %v\n", err)
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
		data.Message = "部屋を削除しました。"

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

// 参加中のRoomの一覧を返す
func (h *RoomHandler) JoinRoomsList(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("User Not Found: %v\n", err)
			http.Error(w, "User Not Found", http.StatusUnauthorized)
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("userUsecase.GetByID error: %v\n", err)
			http.Error(w, fmt.Sprintf("userUsecase.GetByID error: %v", err), http.StatusInternalServerError)
			return
		}

		// joinRoomを格納
		var joinroomslist SentRoomsList
		prooms, err := h.participatingRoomUsecase.GetByUserID(ctx, user.ID)
		if err != nil {
			fmt.Println("データベースとの接続に失敗しました。")
			log.Printf("participatingRoomUsecase.GetByUserID error: %v\n", err)
			http.Error(w, fmt.Sprintf("participatingRoomUsecase.GetByUserID error: %v", err), http.StatusInternalServerError)
			return
		}
		for _, proom := range *prooms {
			joinroomslist.RoomsList = append(joinroomslist.RoomsList, proom.RoomID)
		}

		// jsonに変換
		sentjson, err := json.Marshal(joinroomslist)
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

// Roomの一覧を返す
func (h *RoomHandler) RoomsList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var roomslist SentRoomsList

		// Room一覧取得
		rooms, err := h.roomUsecase.GetAll(ctx)
		if err != nil {
			fmt.Println("データベースとの接続に失敗しました。")
			log.Printf("roomUsecase.GetAll error: %v\n", err)
			http.Error(w, fmt.Sprintf("roomUsecase.GetAll error: %v", err), http.StatusInternalServerError)
			return
		}

		// Roomを格納
		for _, room := range *rooms {
			roomslist.RoomsList = append(roomslist.RoomsList, room.ID)
		}

		// jsonに変換
		sentjson, err := json.Marshal(roomslist)
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
