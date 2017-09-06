package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Coubs struct {
	Coubs []map[string]interface{} `json:"coubs"`
}


func main() {
	bot, err := tgbotapi.NewBotAPI("441447903:AAEWoGH4ETkmwOgk2l7gtP6g2dMGC3QosJM")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if prepare(update.Message.Text) == "сиськи" {

			log.Println("Сиськи!!!!")

			resp, err := http.Get("https://coub.com/api/v2/timeline/tag/boobs?order_by=newest_popular&page=1")

			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)

			coubs := Coubs{}

			err = json.Unmarshal(body, &coubs)

			if err != nil {
				log.Println(err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "https://coub.com/view/" + coubs.Coubs[1]["permalink"].(string))
			bot.Send(msg)

			continue
		}


		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func prepare(msg string) string {
	return strings.ToLower(strings.TrimSpace(msg))
}