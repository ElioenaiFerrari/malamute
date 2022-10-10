include .env

default: dev

.PHONY: dev
.PHONY: deploy

dev: 
	air

deploy:
		go build -ldflags '-s -w' -o bin/main
		./bin/main
