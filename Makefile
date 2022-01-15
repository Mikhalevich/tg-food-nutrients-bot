
BOT_TOKEN=""
GOOGLE_TRANSLATE_CRED_PATH="./tg-translate-creds.json"
USADA_API_KEY="VCJ27UGoY2XTbnw1DLrfRNArIirto9wB9NRvqZIP"

all: build

.PHONY: build
build:
	go build -mod=vendor -o ./bin/bot cmd/bot/main.go

.PHONY: run
run:
	./bin/bot -token=$(BOT_TOKEN) -googlecred=$(GOOGLE_TRANSLATE_CRED_PATH) -usadakey=$(USADA_API_KEY)

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
