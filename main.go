package main

import (
	"errors"
	"os"

	"github.com/Mikhalevich/tg-food-nutrients-bot/bot"
	"github.com/Mikhalevich/tg-food-nutrients-bot/googleconverter"
	"github.com/Mikhalevich/tg-food-nutrients-bot/usadanutrients"
	"github.com/sirupsen/logrus"
)

type params struct {
	BotToken                string
	GoogleTranslateCredPath string
	UsadaApiKey             string
}

func loadParams() (*params, error) {
	tgToken := os.Getenv("FB_TG_BOT_TOKEN")
	if tgToken == "" {
		return nil, errors.New("env for telegram bot token FB_TG_BOT_TOKEN is not specified")
	}

	googleTranslateCredPath := os.Getenv("FB_GOOGLE_TRANSLATE_CRED_PATH")
	if googleTranslateCredPath == "" {
		return nil, errors.New("env for google translate credential path FB_GOOGLE_TRANSLATE_CRED_PATH is not specified")
	}

	usadaApiKey := os.Getenv("FB_USADA_API_KEY")
	if usadaApiKey == "" {
		return nil, errors.New("env for usada api key FB_USADA_API_KEY is not specified")
	}

	return &params{
		BotToken:                tgToken,
		GoogleTranslateCredPath: googleTranslateCredPath,
		UsadaApiKey:             usadaApiKey,
	}, nil
}

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	params, err := loadParams()
	if err != nil {
		logger.Error("load params error: %v", err)
		return
	}

	gc, err := googleconverter.New(params.GoogleTranslateCredPath, logger)
	if err != nil {
		logger.Error("google converter create error: %v", err)
		return
	}

	b, err := bot.NewFoodBot(params.BotToken, gc, usadanutrients.New(params.UsadaApiKey), logger)
	if err != nil {
		logger.Error("create bot error: %v", err)
		return
	}

	if err := b.RunBot(); err != nil {
		logger.Error("run bot error: %v", err)
		return
	}
}
