all: build

.PHONY: build
build:
	go build -mod=vendor -o ./bin/bot cmd/bot/main.go

.PHONY: run
run:
	./bin/bot -config=config/config.yaml

.PHONY: tag
tag:
	docker build -t mikhalevich/tg-food-nutrients-bot:$(TAG) -f ./script/docker/Dockerfile .
	docker push mikhalevich/tg-food-nutrients-bot:$(TAG)

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
