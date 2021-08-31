package googleconverter

import (
	"context"

	"cloud.google.com/go/translate"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type GoogleConverter struct {
	c *translate.Client
	l *logrus.Logger
}

func New(credFilePath string, l *logrus.Logger) (*GoogleConverter, error) {
	client, err := translate.NewClient(context.Background(), option.WithCredentialsFile(credFilePath))
	if err != nil {
		return nil, err
	}

	return &GoogleConverter{
		c: client,
		l: l,
	}, nil
}

func (gc *GoogleConverter) Convert(text string) (string, error) {
	tag, err := gc.detectLanguage(text)
	if err != nil {
		return "", err
	}

	gc.l.Infof("detected %v language: %s\n", tag, text)

	if tag == language.English {
		return text, nil
	}

	t, err := gc.c.Translate(context.Background(), []string{text}, language.English, nil)
	if err != nil {
		return "", err
	}

	return t[0].Text, nil
}

func (gc *GoogleConverter) detectLanguage(text string) (language.Tag, error) {
	detection, err := gc.c.DetectLanguage(context.Background(), []string{text})
	if err != nil {
		return language.Und, err
	}

	if len(detection) <= 0 {
		return language.Und, nil
	}

	if len(detection[0]) != 1 {
		return language.Und, nil
	}

	if detection[0][0].Confidence < 0.5 {
		return language.Und, nil
	}

	return detection[0][0].Language, nil
}
