.PHONY: build docker_build push test

VERSION = $(shell sh ./version.sh)

build:
	GOOS=linux go build

docker_build: build
	docker build -t slofurno/deploy:$(VERSION) -t slofurno/deploy:latest .

push: docker_build
	docker push slofurno/deploy:$(VERSION)
	docker push slofurno/deploy:latest
