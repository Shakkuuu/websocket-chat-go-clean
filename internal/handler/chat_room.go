package handler

import (
	"fmt"

	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"golang.org/x/net/websocket"
)

// クライアントが参加するチャットルーム
type ChatRoom struct {
	ID      string
	Clients map[*websocket.Conn]string
}

// 作成された各ルームを格納
var rooms = make(map[string]*ChatRoom)

func roomInit(rooms *domain.Rooms) {
	for _, room := range *rooms {
		createRoom(room.ID)
	}
}

// RoomMap一覧取得
func getRooms() map[string]*ChatRoom {
	return rooms
}

// 作ったRoomをMapに追加
func createRoom(roomID string) *ChatRoom {
	room := &ChatRoom{
		ID:      roomID,
		Clients: make(map[*websocket.Conn]string),
	}
	rooms[roomID] = room

	return room
}

// RoomMapから削除
func deleteRoom(roomID string) {
	delete(rooms, roomID)
}

// オンラインのユーザー一覧の取得
func getOnlineUsers(roomid string) ([]string, error) {
	var onlineusers []string

	// オンラインのユーザー取得
	onlineusers = append(onlineusers, "匿名")

	// Room一覧取得
	rooms = getRooms()

	// roomがあるか再度確認
	room, exists := rooms[roomid]
	if !exists {
		return onlineusers, fmt.Errorf("this room was not found")
	}

	// Room内のユーザーを格納
	for _, user := range room.Clients {
		onlineusers = append(onlineusers, user)
	}

	return onlineusers, nil
}
