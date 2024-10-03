# websocket-chat-go-clean

## 概要

WebosocketとGo言語を主に使用した、リアルタイム性のあるチャットアプリケーション。

過去に作成した[https://github.com/Shakkuuu/websocket-chat-go](https://github.com/Shakkuuu/websocket-chat-go)をクリーンアーキテクチャで改修しつつ、一部実装を見直しました。

使用方法は[https://github.com/Shakkuuu/websocket-chat-go](https://github.com/Shakkuuu/websocket-chat-go)をご確認ください。

## メモ

- セキュリティ面 XSSとか対策　、バリデーション ok
- チャットの文字列マークダウン形式で送れるようにしたい ok
- チャットの中身変なもの送信できないようにする ok
- ローカルの時はroom.jsのwssをwsにする。
- userとroomのバリデーションはOK。あとはメッセージのバリデーションと対策 ok
