package main

import (
	"context"
)

func main() {
	cfg := LoadConfig()
	InitOAuth(cfg)
	InitDB()

	token := GetToken(context.Background())

	// TODO: 視聴するチャンネルも指定できるようにする
	StartTwitchChatWithToken(cfg.Username, token.AccessToken, "lazvell")

	select {} // 永久ループ
}
