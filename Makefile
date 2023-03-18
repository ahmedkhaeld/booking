BINARY_NAME=jazzApp

build:
	@go mod vendor
	@echo "Building Jazz..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "Jazz built!"

run: build
	@echo "Starting Jazz..."
	@./tmp/${BINARY_NAME} &
	@echo "Jazz started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm tmp/${BINARY_NAME}
	@echo "Cleaned!"

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

stop:
	@echo "Stopping Jazz..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped Jazz!"

restart: stop start