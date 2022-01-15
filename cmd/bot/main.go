package main

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-food-nutrients-bot/internal/bot"
	"github.com/Mikhalevich/tg-food-nutrients-bot/internal/googleconverter"
	"github.com/Mikhalevich/tg-food-nutrients-bot/internal/usadanutrients"
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
		logger.WithError(err).Error("load params error")
		return
	}

	gc, err := googleconverter.New(params.GoogleTranslateCredPath, logger)
	if err != nil {
		logger.WithError(err).Error("google converter create error")
		return
	}

	if err := bot.Run(params.BotToken, gc, usadanutrients.New(params.UsadaApiKey), logger); err != nil {
		logger.WithError(err).Error("run bot error")
		return
	}
}
