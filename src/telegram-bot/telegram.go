package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	telegramAPI "TOTP-telegram/src/API"
)

var token string

func init() {
	f, _ := os.Open("../../.env")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "TG_BOT_TOKEN=") {
			token = strings.TrimPrefix(line, "TG_BOT_TOKEN=")
			break
		}
	}
}

func main() {
	log.Println(token)
	bot := telegramAPI.NewBot("")
	log.Printf("Bot started!! >w<")
	var offset int = 0
	for {
		updates, err := bot.GetUpdates(offset)
		if err != nil {
			log.Println("Error fetching updates:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if updates != nil {
			log.Println("updates:", updates)
		}
		for _, update := range updates {
			log.Println("aa")
			// Update the offset to the latest update_id + 1
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}
			// If the message text is "/start", send "Hello world"
			if strings.TrimSpace(update.Message.Text) == "/start" {
				_, err := bot.SendMessange(update.Message.Chat.ID, "Hello world")
				if err != nil {
					log.Println("Error sending message:", err)
				} else {
					log.Printf("Sent hello world to chat %d", update.Message.Chat.ID)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}
