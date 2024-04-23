build:
	go build

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build

docker-build-and-push:
	@make docker-build
	@make docker-push

docker-build:
	docker build --platform linux/amd64 -t ghcr.io/sikalabs/github-apps-pull-secret-sync .

docker-push:
	docker push ghcr.io/sikalabs/github-apps-pull-secret-sync

docker-build-and-push-arm64:
	@make docker-build-arm64
	@make docker-push-arm64

docker-build-arm64:
	docker build --platform linux/arm64 -t ghcr.io/sikalabs/github-apps-pull-secret-sync:arm64 .

docker-push-arm64:
	docker push ghcr.io/sikalabs/github-apps-pull-secret-sync:arm64
