package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

var (
	tokenFilePath     = ".token.json"
	twitchOAuthConfig *oauth2.Config
)

func InitOAuth(cfg Config) {
	twitchOAuthConfig = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"chat:read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
			TokenURL: "https://id.twitch.tv/oauth2/token",
		},
	}
}

// 認証フロー or ロード
func GetToken(ctx context.Context) *oauth2.Token {
	// 1. ファイルから読み込み
	token, err := loadToken()
	if err == nil {
		// 2. TokenSourceで自動リフレッシュ対応
		ts := twitchOAuthConfig.TokenSource(ctx, token)
		newToken, err := ts.Token()
		if err == nil {
			// トークンが更新されていたら保存
			if newToken.AccessToken != token.AccessToken {
				log.Println("リフレッシュトークンにより新しいトークンを取得しました")
				saveToken(newToken)
			} else {
				log.Println("保存済みトークンを使用します")
			}
			return newToken
		}
		log.Println("トークンリフレッシュに失敗:", err)
	}

	// 3. エラー or 認証未済なら再認証
	log.Println("認証フローを開始します")
	token = StartAuthFlow()
	saveToken(token)
	return token
}

// 認証フローを開始
var tokenChan = make(chan *oauth2.Token)

func StartAuthFlow() *oauth2.Token {
	url := twitchOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Println("以下のURLを開いてTwitchにログインしてください：")
	fmt.Println(url)

	http.HandleFunc("/callback", handleOAuthCallback)
	go http.ListenAndServe(":8080", nil)

	return <-tokenChan
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := twitchOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "トークン交換失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "認証成功！ターミナルに戻ってください。")
	tokenChan <- token
}

// トークン保存
func saveToken(token *oauth2.Token) {
	f, err := os.Create(tokenFilePath)
	if err != nil {
		log.Printf("トークン保存エラー: %v", err)
		return
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

// トークン読み込み
func loadToken() (*oauth2.Token, error) {
	f, err := os.Open(tokenFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)

	// 明示的にトークンの期限をチェック
	if token.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("トークン期限切れ")
	}

	return &token, err
}
