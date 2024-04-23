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
