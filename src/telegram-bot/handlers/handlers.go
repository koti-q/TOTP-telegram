package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	tg "TOTP-telegram/src/API"
	data "TOTP-telegram/src/data"
	totp "TOTP-telegram/src/totp-generator"
)

func HandleHelloWorld(bot tg.BotAPI, chatID int64) {
	_, err := bot.SendMessange(chatID, "Hello world")
	if err != nil {
		log.Println("Error sending message:", err)
	} else {
		log.Printf("Sent hello world to chat %d", chatID)
	}
}

func HandleStart(bot tg.BotAPI, chatID int64) {
	_, err := bot.SendMessange(chatID,
		"Welcome to the TOTP Telegram bot!\n"+
			"/generate {name_secret} {secret} - Generate a TOTP\n"+
			"/send {name} Send a TOTP to the user\n")

	if err != nil {
		log.Println("Error sending message:", err)
	} else {
		log.Printf("Sent to %d", chatID)
	}
}

func HandleGenerateTOTP(bot tg.BotAPI, message tg.Message, key string) {
	m := strings.Split(message.Text, " ")
	if len(m) != 3 {
		bot.SendMessange(message.Chat.ID, "Usage: /generate {name_secret} {secret}")
		return
	}

	_, err := data.GetUser(message.Chat.ID)
	if err != nil {
		data.AddUser(message.Chat.ID)
	}

	name := m[1]
	log.Println(m[2])
	secret, _ := data.EncryptSecret(m[2], key)
	log.Println(secret)
	err = data.SaveSecret(message.Chat.ID, name, secret)
	if err != nil {
		log.Println("Error saving secret:", err)
		bot.SendMessange(message.Chat.ID, "Error saving secret")
		return
	}
	bot.SendMessange(message.Chat.ID, fmt.Sprintf("TOTP generated successfully\n"+
		"To get OTP use: /send %s", name))
}

func HandleSendTOTP(bot tg.BotAPI, message tg.Message, key string) {
	m := strings.Split(message.Text, " ")
	if len(m) != 2 {
		bot.SendMessange(message.Chat.ID, "Usage: /send {name}")
		return
	}
	name := m[1]
	s, err := data.GetSecret(message.Chat.ID, name)
	if err != nil {
		log.Println("Error getting secret:", err)
		bot.SendMessange(message.Chat.ID, "Error getting secret")
		return
	}
	secret, err := data.DecryptSecret(s, key)
	if err != nil {
		log.Println("Error decrypting secret:", err)
		bot.SendMessange(message.Chat.ID, "Error decrypting secret")
		return
	}
	otp := totp.GenerateTOTP(secret, time.Now().Unix())

	bot.SendMessange(message.Chat.ID, fmt.Sprintf("Your TOTP: %d", otp))
}
