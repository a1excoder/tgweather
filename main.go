package main

import (
	"database/sql"
	"fmt"
	"log"

	botLib "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	FlagUa = "🇺🇦"
	FlagRu = "🇷🇺"
	FlagGb = "🇬🇧"
)

func main() {
	data, err := GetConfigFileData("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := GetDataBasePtr("./users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bot, err := botLib.NewBotAPI(data.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	u := botLib.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("bot is listening...")
	for message := range updates {
		if message.Message == nil {
			continue
		}

		if message.Message.CommandWithAt() == "start" {

			if !CheckUser(db, message.Message.From.ID) {
				err = CreateNewUser(db, message.Message.From.ID)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			SendFlags(message, bot)
			continue
		}

		if status, err := SetLanguage(message, db); err != nil {
			log.Println(err)
			continue
		} else if status {
			if message.Message.Text == FlagUa {
				ReplyToMessage(message, bot, fmt.Sprintf("Привіт %s", message.Message.Text))
			} else if message.Message.Text == FlagRu {
				ReplyToMessage(message, bot, fmt.Sprintf("Привет %s", message.Message.Text))
			} else {
				ReplyToMessage(message, bot, fmt.Sprintf("Hello %s", message.Message.Text))
			}

			continue
		}

		if message.Message.CommandWithAt() == "weather" {
			lang, err := GetUserLang(db, message.Message.From.ID)
			if err != nil {
				SendMessage(message, bot, "server error, repeat after 10m")
				continue
			}

			msg, err := GetWeather(message.Message.CommandArguments(), data.OwmToken, lang)
			if err != nil {
				SendMessage(message, bot, "city not found") // server error, repeat after 10m
				log.Println(err)
				continue
			}

			ReplyToMessage(message, bot, msg)
			continue
		}

		_, err = bot.Send(botLib.NewMessage(message.Message.Chat.ID, "command not found"))
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func ReplyToMessage(update botLib.Update, bot *botLib.BotAPI, message string) {
	msg := botLib.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func SendMessage(update botLib.Update, bot *botLib.BotAPI, message string) {
	msg := botLib.NewMessage(update.Message.Chat.ID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func SendFlags(update botLib.Update, bot *botLib.BotAPI) {
	msg := botLib.NewMessage(update.Message.Chat.ID, "select language:")
	languages := []botLib.KeyboardButton{
		{Text: FlagUa},
		{Text: FlagRu},
		{Text: FlagGb},
	}
	msg.ReplyMarkup = botLib.NewReplyKeyboard(languages)

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func SetLanguage(update botLib.Update, db *sql.DB) (bool, error) {
	if update.Message.Text == FlagUa ||
		update.Message.Text == FlagRu ||
		update.Message.Text == FlagGb {
		err := SetUserLang(db, update.Message.From.ID, update.Message.Text)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
