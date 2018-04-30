package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

type TGBOT struct {
	*tb.Bot
	receiver map[int]*tb.Chat
}

func (bot *TGBOT) Init(token, c0, c1 string) (err error) {
	bot.Bot, err = tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	bot.Handle("/id", func(m *tb.Message) {
		bot.Bot.Send(m.Sender, strconv.FormatInt(m.Chat.ID, 10))
	})
	go bot.Start()
	chat0, err0 := bot.Bot.ChatByID(c0)
	chat1, err1 := bot.Bot.ChatByID(c1)
	if err0 != nil || err1 != nil {
		log.Fatal(err0, err1)
		return err0
	}
	bot.receiver = map[int]*tb.Chat{
		0: chat0,
		1: chat1,
	}
	return
}

func (bot *TGBOT) Send(sms SMS) (err error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Println("err", err)
		return
	}

	t := time.Unix(sms.date/1000, 0).In(loc)
	message := fmt.Sprintf("%s\n\n%s\n\nSender: %s", sms.body, t.Format("Recv at 01-02 15:04:05"), sms.address)
	_, err = bot.Bot.Send(bot.receiver[sms.sub_id], message)
	return
}
