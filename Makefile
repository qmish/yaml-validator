# Makefile для yaml-validator
# Использование: make build | test | validate | docker

.PHONY: build test validate docker clean help

BINARY := yaml-validator
GOFLAGS := -v

help:
	@echo "Доступные цели:"
	@echo "  make build    - собрать бинарник"
	@echo "  make test     - запустить тесты"
	@echo "  make validate - проверить тестовые YAML"
	@echo "  make docker   - собрать Docker образ"
	@echo "  make clean    - удалить бинарники"

build:
	go build -o bin/$(BINARY) .

build-cross:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY)-linux .
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY)-mac .
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY).exe .

test:
	go test ./... -v -short

test-coverage:
	go test ./... -coverprofile=coverage.out -short
	go tool cover -html=coverage.out -o coverage.html

validate: build
	./bin/$(BINARY) validate testdata/examples/valid.yaml testdata/examples/k8s-deployment.yaml
	./bin/$(BINARY) validate testdata/examples/docker-compose.yaml -c configs/docker-compose.yaml

docker:
	docker build -t yaml-validator .

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
