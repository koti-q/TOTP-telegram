package main

import (
	"log"
	"os"
	"strings"
	"time"

	tg "TOTP-telegram/src/API"
	data "TOTP-telegram/src/data"

	"github.com/joho/godotenv"

	handler "TOTP-telegram/src/telegram-bot/handlers"
)

var token string
var db_url string

func init() {
	godotenv.Load()

	token = os.Getenv("TG_BOT_TOKEN")
	db_url = os.Getenv("DATABASE_URL")

	if token == "" || db_url == "" {
		log.Fatal("Environment variables TG_BOT_TOKEN and DATABASE_URL must be set")
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
