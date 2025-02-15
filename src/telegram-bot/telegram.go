package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	tg "TOTP-telegram/src/API"
	data "TOTP-telegram/src/data"
	handler "TOTP-telegram/src/telegram-bot/handlers"
)

var token string
var db_url string

func init() {
	f, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "TG_BOT_TOKEN=") {
			token = strings.TrimPrefix(line, "TG_BOT_TOKEN=")
		}
		if strings.HasPrefix(line, "DATABASE=") {
			db_url = strings.TrimPrefix(line, "DATABASE=")
		}
		break
	}
}

var commands = map[string]interface{}{
	"/start":      handler.HandleStart,
	"/helloworld": handler.HandleHelloWorld,
	"/generate":   handler.HandleGenerateTOTP,
	"/send":       handler.HandleSendTOTP,
}

func main() {
	log.Println(token)
	bot := tg.NewBot(token)
	log.Printf("Bot started!! >w<")

	err := data.InitDB(db_url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

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
			// Update the offset to the latest update_id + 1
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}
			// If the message text is "/start", send "Hello world"
			switch strings.Split(update.Message.Text, " ")[0] {
			case "/start":
				handler.HandleStart(*bot, update.Message.Chat.ID)
			case "/helloworld":
				handler.HandleHelloWorld(*bot, update.Message.Chat.ID)
			case "/generate":
				handler.HandleGenerateTOTP(*bot, update.Message)
				offset++
			case "/send":
				handler.HandleSendTOTP(*bot, update.Message)
			default:
				handler.HandleHelloWorld(*bot, update.Message.Chat.ID)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
