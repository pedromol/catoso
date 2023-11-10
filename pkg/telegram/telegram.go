package telegram

import (
	"math/rand"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func (t Telegram) randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (t Telegram) SendPhoto(chatId int64, img []byte) error {
	file := tgbotapi.FileBytes{
		Name:  t.randomString(6) + ".jpg",
		Bytes: img,
	}
	msg := tgbotapi.NewPhoto(chatId, file)
	_, err := t.bot.Send(msg)
	return err
}

func (t Telegram) SendVideo(chatId int64, video []byte) error {
	file := tgbotapi.FileBytes{
		Name:  t.randomString(6) + ".mp4",
		Bytes: video,
	}
	msg := tgbotapi.NewVideo(chatId, file)
	_, err := t.bot.Send(msg)
	return err
}

func NewTelegram(token string) (Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return Telegram{}, err
	}

	return Telegram{
		bot: bot,
	}, nil
}
