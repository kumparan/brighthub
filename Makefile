test:
	richgo test ./... -v --cover

mockgen:
	mockgen -destination=mock/mock_brighthub.go -package=mock github.com/kumparan/brighthub Client

.PHONY: test mockgen