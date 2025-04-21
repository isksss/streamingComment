package main

import (
	"fmt"
	"log"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

func StartTwitchChatWithToken(username, accessToken, channelName string) {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		fmt.Printf("[Twitch: %s] %s: %s\n", channelName, msg.User.Name, msg.Message)
		SaveComment("Twitch", channelName, msg.User.Name, msg.Message)
	})

	client.Join(channelName)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
}
