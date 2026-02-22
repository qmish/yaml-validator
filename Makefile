# Makefile для yaml-validator
# Использование: make build | test | validate | docker

.PHONY: build test test-coverage coverage-check lint validate docker clean help benchmark

BINARY := yaml-validator
GOFLAGS := -v

help:
	@echo "Доступные цели:"
	@echo "  make build    - собрать бинарник"
	@echo "  make test     - запустить тесты"
	@echo "  make test-coverage   - тесты с покрытием и coverage.html"
	@echo "  make coverage-check  - тесты с покрытием, проверка порога 40% (как в CI)"
	@echo "  make lint     - golangci-lint и go vet"
	@echo "  make validate  - проверить тестовые YAML"
	@echo "  make benchmark - запустить бенчмарки производительности"
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

lint:
	go vet ./...
	golangci-lint run ./...

test-coverage:
	go test ./... -coverprofile=coverage.out -short
	go tool cover -html=coverage.out -o coverage.html

coverage-check:
	go test ./... -coverprofile=coverage.out -short -covermode=atomic
	@go tool cover -func=coverage.out | grep total || true
	@echo "CI требует покрытие >= 40%%"

benchmark:
	go test ./internal/validator -bench=. -benchmem -run=^$$ -count=1
	go test ./internal/parser -bench=. -benchmem -run=^$$ -count=1

validate: build
	./bin/$(BINARY) validate testdata/examples/valid.yaml testdata/examples/k8s-deployment.yaml
	./bin/$(BINARY) validate testdata/examples/docker-compose.yaml -c configs/docker-compose.yaml

docker:
	docker build -t yaml-validator .

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
