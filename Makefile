.PHONY: mocks
mocks:
	@go tool mockery --log-level="" && rm -rf ./mocks && mkdir mocks && go tool mockery --log-level=""

.PHONY: fix-imports
fix-imports:
	@go tool goimports -w -local "todo_list" .

.PHONY: test
test:
	@go test -count=1 -covermode=atomic ./...

.PHONY: docker-up
docker-up:
	@docker compose down && docker system prune --volumes --force && docker compose up -d
