.PHONY: lint

node_modules/eslint/bin/eslint.js: package.json yarn.lock
	yarn install

lint: node_modules/eslint/bin/eslint.js
	golangci-lint run
	yarn run eslint web/events.js

build-image: Dockerfile
	docker build -t ghcr.io/petewall/event-logger:dev .
