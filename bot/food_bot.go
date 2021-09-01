package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type NameConverter interface {
	Convert(text string) (string, error)
}

type Nutrienter interface {
	Nutrients(product string, allNutrients bool) (string, error)
}

type FoodBot struct {
	bot *tgbotapi.BotAPI
	c   NameConverter
	n   Nutrienter
	l   *logrus.Logger
}

func NewFoodBot(token string, c NameConverter, n Nutrienter, l *logrus.Logger) (*FoodBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	//bot.Debug = true

	return &FoodBot{
		bot: bot,
		c:   c,
		n:   n,
		l:   l,
	}, nil
}

func (fb *FoodBot) RunBot() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := fb.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			fb.l.Error("invalid message")
			continue
		}

		fb.l.Infof("incoming mesage with chat id = %d, message = %s\n", update.Message.Chat.ID, update.Message.Text)
		go fb.processMessage(update.Message.MessageID, update.Message.Chat.ID, update.Message.Text)
	}

	return nil
}

func (fb *FoodBot) processMessage(messageID int, chatID int64, text string) {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "/start") {
		fb.sendMessage(messageID, chatID, "type product name and you will receive its nutrients")
		return
	}

	allNutrients := false
	if strings.HasPrefix(text, "/all ") {
		text = strings.TrimPrefix(text, "/all ")
		allNutrients = true
	}

	if strings.HasPrefix(text, "/") {
		return
	}

	productName, err := fb.c.Convert(text)

	defer func() {
		if err != nil {
			fb.sendMessage(messageID, chatID, "invalid input")
		}
	}()

	if err != nil {
		fb.l.Errorf("convert error: %v\n", err)
		return
	}

	nutrients, err := fb.n.Nutrients(productName, allNutrients)
	if err != nil {
		fb.l.Errorf("nutrients error: %v\n", err)
		return
	}

	fb.sendMessage(messageID, chatID, nutrients)
}

func (fb *FoodBot) sendMessage(messageID int, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID
	fb.bot.Send(msg)
}
