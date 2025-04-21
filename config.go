package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientID     string
	ClientSecret string
	Username     string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("環境変数の読み込みに失敗しました: %v", err)
	}

	return Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		Username:     os.Getenv("TWITCH_USERNAME"),
	}
}
