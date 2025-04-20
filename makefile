.PHONY: run-be run-fe install-fe deps-be build-be clean docker-build docker-up docker-down test-be generate-swagger

# Backend
run-be:
	cd be/cmd/api && go run main.go

deps-be:
	cd be && go mod tidy

build-be:
	cd be && go build -o bin/api cmd/api/main.go

# Frontend
install-fe:
	cd fe && npm install

run-fe:
	cd fe && npm run dev

build-fe:
	cd fe && npm run build

test-be:
	cd be && go test ./...


# Combined
run: run-be run-fe

# Docker
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Swagger
generate-swagger:
	cd be && swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --parseDepth 2 --generatedTime=false --propertyStrategy=camelcase

# Clean
clean:
	rm -rf be/bin
	rm -rf fe/.next
	rm -rf fe/out
	rm -rf fe/node_modules
	docker-compose down -v