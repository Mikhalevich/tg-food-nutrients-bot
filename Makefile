
BOT_TOKEN=""
GOOGLE_TRANSLATE_CRED_PATH="./tg-translate-creds.json"
USADA_API_KEY="VCJ27UGoY2XTbnw1DLrfRNArIirto9wB9NRvqZIP"

all: build

.PHONY: build
build:
	go build -mod=vendor -o ./bin/bot cmd/bot/main.go

.PHONY: run
run:
	FB_TG_BOT_TOKEN=$(BOT_TOKEN) FB_GOOGLE_TRANSLATE_CRED_PATH=$(GOOGLE_TRANSLATE_CRED_PATH) FB_USADA_API_KEY=$(USADA_API_KEY) ./bin/bot

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
