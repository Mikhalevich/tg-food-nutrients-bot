package main

import (
	"flag"
	"os"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-food-nutrients-bot/internal/bot"
	"github.com/Mikhalevich/tg-food-nutrients-bot/internal/nameconverter"
	"github.com/Mikhalevich/tg-food-nutrients-bot/internal/usadanutrients"
)

type config struct {
	BotToken                string `yaml:"bot_token" required:"true"`
	GoogleTranslateCredPath string `yaml:"google_transtale_cred_path"`
	UsadaApiKey             string `yaml:"usada_api_key" required:"true"`
}

func isDebugEnabled() bool {
	if debug := os.Getenv("FB_DEBUG"); debug != "" {
		return true
	}
	return false
}

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	configFile := flag.String("config", "config/config.yaml", "telegram bot config file")
	flag.Parse()

	var cfg config
	if err := configor.Load(&cfg, *configFile); err != nil {
		logger.WithError(err).Error("failed to load bot config")
		return
	}

	converter, err := nameconverter.MakeConverter(cfg.GoogleTranslateCredPath, logger)
	if err != nil {
		logger.WithError(err).Error("google converter create error")
		return
	}

	if err := bot.Run(cfg.BotToken, converter, usadanutrients.New(cfg.UsadaApiKey), logger, isDebugEnabled()); err != nil {
		logger.WithError(err).Error("run bot error")
		return
	}
}
