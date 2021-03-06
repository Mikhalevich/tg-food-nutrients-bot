package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type nameConverter interface {
	Convert(text string) (string, error)
}

type nutrienter interface {
	Nutrients(product string, allNutrients bool) (string, error)
}

type foodBot struct {
	bot *tgbotapi.BotAPI
	c   nameConverter
	n   nutrienter
	l   *logrus.Logger
}

func async(fn func() error) <-chan error {
	ch := make(chan error)

	go func() {
		ch <- fn()
	}()

	return ch
}

func Run(token string, c nameConverter, n nutrienter, l *logrus.Logger, debugEnabled bool) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("create bot api: %w", err)
	}
	bot.Debug = debugEnabled

	fb := &foodBot{
		bot: bot,
		c:   c,
		n:   n,
		l:   l,
	}

	done := async(func() error {
		l.Info("bot running...")
		return fb.processUpdates()
	})

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-terminate:
			signal.Stop(terminate)
			l.Info("stopping bot...")
			bot.StopReceivingUpdates()
		case err = <-done:
			break loop
		}
	}

	l.Info("bot stopped...")

	return err
}

func (fb *foodBot) processUpdates() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	updates := fb.bot.GetUpdatesChan(u)

	var wg sync.WaitGroup

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			fb.l.Error("invalid message")
			continue
		}

		fb.l.WithFields(logrus.Fields{
			"chat_id": update.Message.Chat.ID,
			"message": update.Message.Text,
		}).Info("incoming message")

		wg.Add(1)

		go func(u tgbotapi.Update) {
			defer wg.Done()
			fb.processMessage(u.Message.MessageID, u.Message.Chat.ID, u.Message.Text)
		}(update)
	}

	wg.Wait()

	return nil
}

func (fb *foodBot) processMessage(messageID int, chatID int64, text string) {
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

	nutrients, err := fb.nutrients(text, allNutrients)
	if err != nil {
		fb.l.WithError(err).Error("nutrients error")
		fb.sendMessage(messageID, chatID, "invalid input")
		return
	}

	fb.sendMessage(messageID, chatID, nutrients)
}

func (fb *foodBot) nutrients(food string, all bool) (string, error) {
	food, err := fb.c.Convert(food)
	if err != nil {
		return "", fmt.Errorf("food convert error: %w", err)
	}

	nutrients, err := fb.n.Nutrients(food, all)
	if err != nil {
		return "", fmt.Errorf("nutrients error: %w", err)
	}

	return nutrients, nil
}

func (fb *foodBot) sendMessage(messageID int, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = messageID
	fb.bot.Send(msg)
}
