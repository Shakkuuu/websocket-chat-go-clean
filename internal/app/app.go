package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shakkuuu/websocket-chat-go-clean/config"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/domain"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/handler"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/repository"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/usecase"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/httpserver"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/postgres"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/session"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/timefmt"
	"github.com/Shakkuuu/websocket-chat-go-clean/pkg/ulid"
	"golang.org/x/net/websocket"
)

var accesslogfile *os.File

func Run(cfg *config.Config, accessfile, chatLogFile *os.File) {
	accesslogfile = accessfile

	pg, err := postgres.New(cfg.ProtocolDB, cfg.UserNameDB, cfg.UserPassDB, cfg.NameDB, cfg.PortDB)
	if err != nil {
		log.Fatal(fmt.Errorf("app.Run.postgres.New: %w", err))
	}
	defer pg.Close()
	pg.Db.AutoMigrate(&domain.User{})
	pg.Db.AutoMigrate(&domain.ParticipatingRoom{})
	pg.Db.AutoMigrate(&domain.Room{})
	insertTokumei(pg)

	newSession := session.New()
	mux := http.NewServeMux()

	userRepo := repository.NewUserRepo(pg)
	participatingRoomRepo := repository.NewParticipatingRoomRepo(pg)
	roomRepo := repository.NewRoomRepo(pg)
	userUsecase := usecase.NewUserUsecase(userRepo)
	participatingRoomUsecase := usecase.NewParticipatingRoomUsecase(participatingRoomRepo)
	roomUsecase := usecase.NewRoomUsecase(roomRepo)

	// User
	userHandler := handler.NewUserHandler(userUsecase, participatingRoomUsecase, roomUsecase, newSession)
	mux.Handle("/usermenu", loggingMiddleware(http.HandlerFunc(userHandler.Menu)))                 // usermenuページ
	mux.Handle("/login", loggingMiddleware(http.HandlerFunc(userHandler.Login)))                   // ログインページ
	mux.Handle("/signup", loggingMiddleware(http.HandlerFunc(userHandler.Signup)))                 // サインアップページ
	mux.Handle("/logout", loggingMiddleware(http.HandlerFunc(userHandler.Logout)))                 // ログアウト処理
	mux.Handle("/deleteuser", loggingMiddleware(http.HandlerFunc(userHandler.Delete)))             // User削除
	mux.Handle("/changepassword", loggingMiddleware(http.HandlerFunc(userHandler.ChangePassword))) // パスワード更新
	mux.Handle("/username", loggingMiddleware(http.HandlerFunc(userHandler.GetUserName)))          // 自身のユーザー名取得

	// Room
	roomHandler := handler.NewRoomHandler(userUsecase, participatingRoomUsecase, roomUsecase, newSession)
	mux.Handle("/", loggingMiddleware(http.HandlerFunc(roomHandler.Top)))                    // roomtopページ
	mux.Handle("/room", loggingMiddleware(http.HandlerFunc(roomHandler.Room)))               // Room内のページ
	mux.Handle("/deleteroom", loggingMiddleware(http.HandlerFunc(roomHandler.Delete)))       // Room削除
	mux.Handle("/rooms", loggingMiddleware(http.HandlerFunc(roomHandler.RoomsList)))         // Room一覧取得
	mux.Handle("/joinrooms", loggingMiddleware(http.HandlerFunc(roomHandler.JoinRoomsList))) // 参加中のRoom一覧取得

	// websocket
	websocketHandler := handler.NewWebsocketHandler(userUsecase, participatingRoomUsecase, roomUsecase, newSession, chatLogFile)
	mux.Handle("/ws", websocket.Handler(websocketHandler.HandleConnection)) // メッセージWebsocket用
	go websocketHandler.HandleMessages()                                    // goroutineとチャネルで常にメッセージを待つ

	// static
	staticFileDirectory := http.Dir("./static")
	staticFileServer := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	mux.Handle("/static/", staticFileServer)

	httpServer := httpserver.New(mux, httpserver.Port(cfg.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Fatal(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		accesslog := fmt.Sprintf("%s: [%s] %s %s %s\n", timefmt.TimeToStr(start), r.Method, r.RemoteAddr, r.URL, time.Since(start))
		fmt.Print(accesslog)
		fmt.Fprint(accesslogfile, accesslog)
	})
}

// 匿名ユーザーを初期に追加
func insertTokumei(pg *postgres.Postgres) {
	pg.Db.Where("name = ?", "匿名").Delete(&domain.User{})

	tokumei := domain.User{
		ID:       ulid.NewULID(),
		Name:     "匿名",
		Password: "tokumei",
	}

	err := pg.Db.Create(&tokumei).Error
	if err != nil {
		log.Printf("db.Create tokumei error: %v\n", err)
	}
	log.Println("匿名ユーザーが登録されました。")
}
