package main

import (
	"errors"
	"flag"
	"fmt"
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
	Debug                   bool
}

func (p params) validate() error {
	if p.BotToken == "" {
		return errors.New("telegram bot token is required")
	}

	if p.GoogleTranslateCredPath == "" {
		return errors.New("google translate credentials is requred")
	}

	if p.UsadaApiKey == "" {
		return errors.New("usada api key is required")
	}

	return nil
}

func loadParams() (*params, error) {
	var p params
	flag.StringVar(&p.BotToken, "token", "", "telegram bot token")
	flag.StringVar(&p.GoogleTranslateCredPath, "googlecred", "", "path for google translate credentials file")
	flag.StringVar(&p.UsadaApiKey, "usadakey", "", "usada api key https://fdc.nal.usda.gov/api-key-signup.html")

	flag.Parse()

	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("params validation: %w", err)
	}

	if debug := os.Getenv("FB_DEBUG"); debug != "" {
		p.Debug = true
	}

	return &p, nil
}

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	p, err := loadParams()
	if err != nil {
		logger.WithError(err).Error("load params error")
		return
	}

	gc, err := googleconverter.New(p.GoogleTranslateCredPath, logger)
	if err != nil {
		logger.WithError(err).Error("google converter create error")
		return
	}

	if err := bot.Run(p.BotToken, gc, usadanutrients.New(p.UsadaApiKey), logger, p.Debug); err != nil {
		logger.WithError(err).Error("run bot error")
		return
	}
}
