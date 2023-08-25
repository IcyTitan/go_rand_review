package telegramBot

import (
	"database/sql"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

var token = ""

func InitBot() {
	startBot()
}

func SetToken(telegramToken string) {
	token = telegramToken
}

func startBot() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 15

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := strings.Fields(update.Message.Text)

		switch message[0] {
		case "/review":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Кто будет делать ревью?")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("!Front"),
					tgbotapi.NewKeyboardButton("!Back"),
				),
			)
			bot.Send(msg)
		case "!Back":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, randBack())
			bot.Send(msg)

		case "!Front":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, randFront())
			bot.Send(msg)

		case "Add_Front":
			if len(message) == 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, addFront(message[1]))
				bot.Send(msg)
			}

		case "Add_Back":
			if len(message) == 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, addBack(message[1]))
				bot.Send(msg)
			}
		}

	}
}

type users struct {
	id   int
	name string
}

func randFront() string {
	db, err := sql.Open("sqlite3", "/var/bots/store.db")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select name from front ORDER BY RANDOM () LIMIT 1")
	user := "никого нет :("

	for rows.Next() {
		rows.Scan(&user)
	}

	if err != nil {
		return user
	}

	defer db.Close()

	return user
}

func randBack() string {
	db, err := sql.Open("sqlite3", "/var/bots/store.db")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select name from back ORDER BY RANDOM () LIMIT 1")
	user := "никого нет :("

	for rows.Next() {
		rows.Scan(&user)
	}

	if err != nil {
		return user
	}

	defer db.Close()

	return user
}

func addBack(nickname string) string {

	db, err := sql.Open("sqlite3", "/var/bots/store.db")
	if err != nil {
		panic(err)
	}

	db.Exec("insert into back (name) values ($1)", nickname)

	defer db.Close()

	return nickname
}

func addFront(nickname string) string {

	db, err := sql.Open("sqlite3", "/var/bots/store.db")
	if err != nil {
		panic(err)
	}

	db.Exec("insert into front (name) values ($1)", nickname)

	defer db.Close()

	return nickname
}
