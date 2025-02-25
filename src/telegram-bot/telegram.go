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
var encryptionKey string

type CommandHandler func(bot *tg.BotAPI, msg tg.Message)

func init() {
	godotenv.Load()

	token = os.Getenv("TG_BOT_TOKEN")
	db_url = os.Getenv("DATABASE_URL")
	encryptionKey = os.Getenv("ENCRYPTION_KEY")

	if token == "" || db_url == "" {
		log.Fatal("Environment variables TG_BOT_TOKEN and DATABASE_URL must be set")
	}
}

func main() {
	bot := tg.NewBot(token)
	log.Printf("Bot started!! >w<")

	err := data.InitDB(db_url)
	if err != nil {
		for err != nil {
			log.Println("Error connecting to database:", err)
			time.Sleep(2 * time.Second)
			err = data.InitDB(db_url)
		}
	}
	log.Println("Connected to database")

	commands := map[string]CommandHandler{
		"/start": func(bot *tg.BotAPI, msg tg.Message) {
			handler.HandleStart(*bot, msg.Chat.ID)
		},
		"/helloworld": func(bot *tg.BotAPI, msg tg.Message) {
			handler.HandleHelloWorld(*bot, msg.Chat.ID)
		},
		"/generate": func(bot *tg.BotAPI, msg tg.Message) {
			handler.HandleGenerateTOTP(*bot, msg, encryptionKey)
		},
		"/send": func(bot *tg.BotAPI, msg tg.Message) {
			handler.HandleSendTOTP(*bot, msg, encryptionKey)
		},
		"/s": func(bot *tg.BotAPI, msg tg.Message) {
			handler.HandleSendTOTP(*bot, msg, encryptionKey)
		},
	}

	offset := 0
	for {
		updates, err := bot.GetUpdates(offset)
		if err != nil {
			log.Println("Error fetching updates:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if updates != nil {
			log.Println("Received updates:", updates)
		}
		for _, update := range updates {
			// Update the offset so we don't reprocess the same update.
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}
			// Extract the command from the message.
			command := strings.Split(update.Message.Text, " ")[0]
			if cmdHandler, ok := commands[command]; ok {
				cmdHandler(bot, update.Message)
			} else {
				// Fallback to a default handler if the command is not recognized.
				handler.HandleHelloWorld(*bot, update.Message.Chat.ID)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
