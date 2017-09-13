package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type User struct {
	ID       int
	Username string
	ChatID   int
	Coubs    *Coubs
}

var Users map[int]*User
var Data map[string]string

func main() {

	Users = make(map[int]*User)

	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	db, err := ioutil.ReadFile("db.json")
	if err != nil {
		log.Panic(err)
	}

	Data = make(map[string]string)
	err = json.Unmarshal(db, &Data)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// chatConfig = ChatConfig{ChatID:}
	// chat, _ := bot.GetChat(chatConfig) // ChatID https://t.me/joinchat/B3GjBA45483QcvfLTR_JNg

	for update := range updates {
		if update.Message == nil {
			continue
		}
		user, ok := Users[update.Message.From.ID]

		if !ok {
			log.Println("New user!!!")
			user = &User{ID: update.Message.From.ID, ChatID: int(update.Message.Chat.ID), Username: update.Message.From.UserName, Coubs: &Coubs{}}
			Users[user.ID] = user
		}
		if prepare(update.Message.Text) == "/help" {
			text := `
			Хочешь сиськи пиши "покажи сиськи"
			Хочешь животных пиши "покажи животных"
			Хочешь смеятся пиши "покажи забавное"

			написал что хотел, для продолжения пиши "еще" 
			или просто "\"

			для обучения пиши: !вопрос == ответ
			`
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
			continue
		}

		if prepare(update.Message.Text) == "покажи сиськи" {
			permalink := ""
			if user.Coubs.Tag == "boobs" {
				permalink = user.Coubs.GetNext()
			} else {
				permalink = user.Coubs.GetBoobs()
			}
			log.Println(permalink)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, permalink)
			bot.Send(msg)
			continue
		}

		if prepare(update.Message.Text) == "покажи животных" {
			permalink := ""
			if user.Coubs.Tag == "animals-pets" {
				permalink = user.Coubs.GetNext()
			} else {
				permalink = user.Coubs.GetAnimal()
			}
			log.Println(permalink)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, permalink)
			bot.Send(msg)
			continue
		}

		if prepare(update.Message.Text) == "покажи забавное" {
			permalink := ""
			if user.Coubs.Tag == "funny" {
				permalink = user.Coubs.GetNext()
			} else {
				permalink = user.Coubs.GetFunny()
			}
			log.Println(permalink)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, permalink)
			bot.Send(msg)
			continue
		}

		if prepare(update.Message.Text) == "еще" || prepare(update.Message.Text) == "\\" {
			permalink := user.Coubs.GetNext()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, permalink)
			bot.Send(msg)
			continue
		}

		ansver, ok := Data[prepare(update.Message.Text)]
		if ok {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansver)
			bot.Send(msg)
			continue
		}

		if strings.Index(prepare(update.Message.Text), "!") == 0 {
			text := strings.Split(strings.TrimLeft(strings.TrimSpace(update.Message.Text), "!"), "==")
			if len(text) == 2 {
				Data[prepare(text[0])] = strings.TrimSpace(text[1])
				jsonData, err := json.Marshal(Data)
				if err != nil {
					log.Println(err)
				}
				err = ioutil.WriteFile("db.json", jsonData, 0755)
				if err != nil {
					log.Println(err)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Запомнил!")
					bot.Send(msg)
					continue
				}
			}
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)
	}
}

func prepare(msg string) string {
	return strings.ToLower(strings.TrimSpace(msg))
}
