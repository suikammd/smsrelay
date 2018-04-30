package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	var db_path, last_path, token string
	var chat0, chat1 string
	var sleep_time int
	flag.StringVar(&db_path, "db", "/tmp/db.db", "database dump path")
	flag.StringVar(&last_path, "last", "last", "last save path")
	flag.StringVar(&token, "token", "", "telegram token")
	flag.StringVar(&chat0, "chat0", "", "chatid 0")
	flag.StringVar(&chat1, "chat1", "", "chatid 1")
	flag.IntVar(&sleep_time, "sleep", 10, "sleep time in second")
	flag.Parse()

	bot := TGBOT{}
	bot.Init(token, chat0, chat1)
	db := Database{last_path: last_path, db_path: db_path}
	db.Init()
	for {
		smses, err := db.Read()
		if err == nil {
			var last int64
			for _, sms := range smses {
				if bot.Send(sms) != nil {
					log.Println(err)
				} else {
					last = sms.date
				}
			}
			db.Save(last)
		} else {
			log.Println(err)
		}
		time.Sleep(time.Duration(sleep_time) * time.Second)
	}
}
