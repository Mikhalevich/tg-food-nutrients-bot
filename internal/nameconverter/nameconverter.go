package nameconverter

import (
	"github.com/sirupsen/logrus"

	"github.com/Mikhalevich/tg-food-nutrients-bot/internal/nameconverter/googleconverter"
	"github.com/Mikhalevich/tg-food-nutrients-bot/internal/nameconverter/stub"
)

type Converter interface {
	Convert(text string) (string, error)
}

func MakeConverter(credFilePath string, logger *logrus.Logger) (Converter, error) {
	if credFilePath == "" {
		logger.Info("using stub converter")
		return stub.New(), nil
	}

	logger.Info("using google converter")
	return googleconverter.New(credFilePath, logger)
}
