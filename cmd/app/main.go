package main

import (
	"io"
	"log"
	"os"

	"github.com/Shakkuuu/websocket-chat-go-clean/config"
	"github.com/Shakkuuu/websocket-chat-go-clean/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// アクセスログ出力用ファイル読み込み
	f, err := os.OpenFile("logs/access.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("logfile os.OpenFile error:%v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	// エラーログ出力用ファイル読み込み
	errorfile, err := os.OpenFile("logs/error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("logfile os.OpenFile error:%v\n", err)
		os.Exit(1)
	}
	defer errorfile.Close()
	// chatログ出力用ファイル読み込み
	chatfile, err := os.OpenFile("logs/chat.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("logfile os.OpenFile error:%v\n", err)
		os.Exit(1)
	}
	defer chatfile.Close()

	// ログの先頭に日付時刻とファイル名、行数を表示するように設定
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// エラーログの出力先をファイルに指定
	log.SetOutput(io.MultiWriter(os.Stderr, errorfile))

	// Run
	app.Run(cfg, f, chatfile)
}
