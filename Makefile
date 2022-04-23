all: build

.PHONY: build
build:
	go build -mod=vendor -o ./bin/bot cmd/bot/main.go

.PHONY: run
run:
	./bin/bot -config=config/config.yaml

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
