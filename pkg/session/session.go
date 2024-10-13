package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const SESSION_NAME string = "Shakkuuu-websocket-chat-go"

type Sessions struct {
	Session *sessions.Session
	Store   *sessions.CookieStore
}

// セッションの初期化
func New() *Sessions {
	s := &Sessions{}

	randBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, randBytes)
	if err != nil {
		panic(err)
	}
	sessionKey := base64.RawURLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes)

	s.Store = sessions.NewCookieStore([]byte(sessionKey))
	s.Session = sessions.NewSession(s.Store, SESSION_NAME)

	return s
}

func (s *Sessions) GetUserData(r *http.Request) (string, string, error) {
	// セッション読み取り
	session, err := s.Store.Get(r, SESSION_NAME)
	if err != nil {
		log.Printf("store.Get error: %v\n", err)
		return "", "", err
	}

	id := session.Values["id"]
	if id == nil {
		fmt.Println("セッションなし")
		err = fmt.Errorf("セッションなし")
		return "", "", err
	}
	i := id.(string)

	username := session.Values["username"]
	if username == nil {
		fmt.Println("セッションなし")
		err = fmt.Errorf("セッションなし")
		return "", "", err
	}
	un := username.(string)

	return i, un, nil
}

func (s *Sessions) Set(r *http.Request, w http.ResponseWriter, id, username string) error {
	s.Session, _ = s.Store.Get(r, SESSION_NAME)
	s.Session.Values["id"] = id
	s.Session.Values["username"] = username
	return s.Session.Save(r, w)
}

func (s *Sessions) Delete(r *http.Request, w http.ResponseWriter) error {
	// セッション削除
	s.Session.Options.MaxAge = -1
	return s.Session.Save(r, w)
}
