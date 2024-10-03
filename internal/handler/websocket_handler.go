package handler

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/usecase"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/session"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/timefmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/net/websocket"
	"gorm.io/gorm"
)

type WebsocketHandler struct {
	userUsecase              usecase.UserUsecase
	participatingRoomUsecase usecase.ParticipatingRoomUsecase
	roomUsecase              usecase.RoomUsecase
	templates                *template.Template
	session                  *session.Sessions
	chatLogFile              *os.File
}

func NewWebsocketHandler(
	usecase usecase.UserUsecase,
	participatingRoomUsecase usecase.ParticipatingRoomUsecase,
	roomUsecase usecase.RoomUsecase,
	session *session.Sessions,
	chatLogFile *os.File,
) *WebsocketHandler {
	templates := template.Must(template.ParseGlob("internal/handler/templates/*.html"))
	return &WebsocketHandler{
		userUsecase:              usecase,
		participatingRoomUsecase: participatingRoomUsecase,
		roomUsecase:              roomUsecase,
		templates:                templates,
		session:                  session,
		chatLogFile:              chatLogFile,
	}
}

var sentmessage = make(chan Message) // 各クライアントに送信するためのメッセージのチャネル

// WebsocketでRoom参加後のコネクション確立
func (h *WebsocketHandler) HandleConnection(ws *websocket.Conn) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 24*time.Hour)
	defer cancel()

	defer ws.Close()

	// セッション読み取り
	userID, userName, err := h.session.GetUserData(ws.Request())
	if err != nil {
		log.Printf("session.GetUserData error: %v\n", err)
		return
	}

	// クライアントから参加する部屋が指定されたメッセージ受信
	var msg Message
	err = websocket.JSON.Receive(ws, &msg)
	if err != nil {
		log.Printf("Receive room ID error:%v\n", err)
		return
	}

	// Room一覧取得
	rooms := getRooms()

	// 部屋が存在しているかどうか
	room, exists := rooms[msg.RoomID]
	if !exists {
		log.Printf("This room was not found\n")
		return
	}

	var notExists bool = false
	_, err = h.participatingRoomUsecase.GetByUserIDAndRoomID(ctx, userID, room.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		notExists = true
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("participatingRoomUsecase.GetByUserIDAndRoomID error: %v\n", err)
		return
	}

	if notExists {
		// 参加中のルーム一覧に参加者として追加
		proom := domain.ParticipatingRoom{
			RoomID:   room.ID,
			IsMaster: false,
			UserID:   userID,
		}
		err := h.participatingRoomUsecase.Create(ctx, &proom)
		if err != nil {
			log.Printf("participatingRoomUsecase.Create: %v\n", err)
			return
		}
	}

	// Roomに参加
	room.Clients[ws] = userName

	// 参加しているユーザー一覧とオンラインのユーザー一覧の取得
	allusersChan := make(chan interface{})
	onlineusersChan := make(chan interface{})
	var allusers []string
	var onlineusers []string
	go func() {
		users, err := h.participatingRoomUsecase.GetUsersByRoomID(ctx, room.ID)
		if err != nil {
			err = fmt.Errorf("participatingRoomUsecase.GetUsersByRoomID error: %v", err)
			allusersChan <- err
			return
		}
		var aus []string
		aus = append(aus, "匿名")
		for _, user := range *users {
			aus = append(aus, user.Name)
		}
		allusersChan <- aus
	}()
	go func() {
		ous, err := getOnlineUsers(room.ID)
		if err != nil {
			err = fmt.Errorf("getOnlineUsers error: %v", err)
			onlineusersChan <- err
			return
		}
		onlineusersChan <- ous
	}()

	auc := <-allusersChan
	ouc := <-onlineusersChan
	switch auctype := auc.(type) {
	case error:
		log.Println(auctype)
		return
	case []string:
		allusers = auctype
	}

	switch ouctype := ouc.(type) {
	case error:
		log.Println(ouctype)
		return
	case []string:
		onlineusers = ouctype
	}

	// Roomに参加したことをそのRoomのクライアントにブロードキャスト
	entermsg := Message{RoomID: room.ID, Message: userName + "が入室しました", Name: "Server", ToName: "", AllUsers: allusers, OnlineUsers: onlineusers}
	sentmessage <- entermsg

	// サーバ側からクライアントにWellcomeメッセージを送信
	err = websocket.JSON.Send(ws, Message{RoomID: room.ID, Message: "ルーム" + room.ID + "へようこそ", Name: "Server", ToName: msg.Name, AllUsers: nil, OnlineUsers: nil})
	if err != nil {
		log.Printf("server wellcome Send error:%v\n", err)
	}

	// クライアントからメッセージが来るまで受信待ちする
	for {
		// クライアントからのメッセージを受信
		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			if err.Error() == "EOF" { // Roomを退出したことを示すメッセージが来たら
				log.Printf("EOF error:%v\n", err)
				delete(room.Clients, ws) // Roomからそのクライアントを削除

				// 参加しているユーザー一覧とオンラインのユーザー一覧の取得
				allusersChan := make(chan interface{})
				onlineusersChan := make(chan interface{})
				var allusers []string
				var onlineusers []string
				go func() {
					users, err := h.participatingRoomUsecase.GetUsersByRoomID(ctx, room.ID)
					if err != nil {
						err = fmt.Errorf("participatingRoomUsecase.GetUsersByRoomID error: %v", err)
						allusersChan <- err
						return
					}
					var aus []string
					aus = append(aus, "匿名")
					for _, user := range *users {
						aus = append(aus, user.Name)
					}
					allusersChan <- aus
				}()
				go func() {
					ous, err := getOnlineUsers(room.ID)
					if err != nil {
						err = fmt.Errorf("getOnlineUsers error: %v", err)
						onlineusersChan <- err
						return
					}
					onlineusersChan <- ous
				}()

				auc := <-allusersChan
				ouc := <-onlineusersChan
				switch auctype := auc.(type) {
				case error:
					log.Println(auctype)
					return
				case []string:
					allusers = auctype
				}

				switch ouctype := ouc.(type) {
				case error:
					log.Println(ouctype)
					return
				case []string:
					onlineusers = ouctype
				}

				// そのクライアントがRoomから退出したことをそのRoomにブロードキャスト
				exitmsg := Message{RoomID: msg.RoomID, Message: msg.Name + "が退出しました", Name: "Server", ToName: "", AllUsers: allusers, OnlineUsers: onlineusers}
				sentmessage <- exitmsg
				break
			}
			log.Printf("Receive error:%v\n", err)
		}

		htmlmsg := blackfriday.Run([]byte(msg.Message))
		policy := bluemonday.UGCPolicy()
		sanitizedHTML := policy.SanitizeBytes(htmlmsg)
		msg.Message = string(sanitizedHTML)

		// goroutineでチャネルを待っているとこへメッセージを渡す
		sentmessage <- msg
	}
}

// goroutineでメッセージのチャネルが来るまで待ち、Roomにメッセージを送信する
func (h *WebsocketHandler) HandleMessages() {
	for {
		// sentmessageチャネルからメッセージを受け取る
		msg := <-sentmessage

		// Room一覧取得
		rooms := getRooms()

		// 部屋が存在しているかどうか
		room, exists := rooms[msg.RoomID]
		if !exists {
			continue
		}

		// チャットログを出力と保存 日時、サーバー名、ユーザー名、宛先、メッセージ
		replaceNlMsg := strings.ReplaceAll(msg.Message, "\n", " ") // 改行があるとログが改行されてしまうため、改行を削除
		chatlog := fmt.Sprintf("%s: [S%s] From(%s) To (%s) Msg(%s)\n", timefmt.TimeToStr(time.Now()), msg.RoomID, msg.Name, msg.ToName, replaceNlMsg)
		fmt.Print(chatlog)
		fmt.Fprint(h.chatLogFile, chatlog)

		if msg.ToName != "" {
			// 接続中のクライアントにメッセージを送る
			for client, name := range room.Clients {
				if msg.ToName == name || msg.Name == name {
					// メッセージを返信する
					policy := bluemonday.UGCPolicy()
					msg.ToName = policy.Sanitize(msg.ToName)
					err := websocket.JSON.Send(client, Message{RoomID: room.ID, Message: msg.Message, Name: msg.Name, ToName: msg.ToName, AllUsers: msg.AllUsers, OnlineUsers: msg.OnlineUsers})
					if err != nil {
						log.Printf("Send error:%v\n", err)
					}
				}
			}
		} else {
			// 接続中のクライアントにメッセージを送る
			for client := range room.Clients {
				// メッセージを返信する
				err := websocket.JSON.Send(client, Message{RoomID: room.ID, Message: msg.Message, Name: msg.Name, ToName: "", AllUsers: msg.AllUsers, OnlineUsers: msg.OnlineUsers})
				if err != nil {
					log.Printf("Send error:%v\n", err)
				}
			}
		}
	}
}
