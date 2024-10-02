FROM golang:1.23.1

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod tidy

COPY ./ ./

EXPOSE 8080

RUN go build -o /shakku-websocket-chat ./cmd/app

# 起動コマンド
CMD ["/shakku-websocket-chat"]
