.PHONY: test
test:
	docker run --rm \
	-v $(PWD):/go/src/github.com/search2d/go-pixiv \
	-w /go/src/github.com/search2d/go-pixiv \
	golang:1.9.0-alpine3.6 \
	go vet ./... && go test -v -cover ./...

.PHONY: integration
integration:
	docker run --rm \
	-v $(PWD):/go/src/github.com/search2d/go-pixiv \
	-w /go/src/github.com/search2d/go-pixiv \
	--env-file .env \
	golang:1.9.0-alpine3.6 \
	go test -v -cover -tags 'integration' ./...