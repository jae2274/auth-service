BINARY_NAME=myapp
CODE_DIR=./usersvc

include test.env

## build: Build binary
build:
	@echo "Building..."
	@go build -ldflags="-s -w" -o ${BINARY_NAME} ${CODE_DIR}
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
	@GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID} GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET} GOOGLE_REDIRECT_URL=${GOOGLE_REDIRECT_URL} SECRET_KEY=${SECRET_KEY} API_PORT=${API_PORT} DB_HOST=${DB_HOST} DB_PORT=${DB_PORT} DB_NAME=${DB_NAME} DB_USERNAME=${DB_USERNAME} DB_PASSWORD=${DB_PASSWORD} ./${BINARY_NAME} 
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

proto:
	@export PATH="$PATH:$(go env GOPATH)/bin"
	@protoc careerhub/provider/provider_grpc/*.proto  --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative  --proto_path=.

## test: runs all tests
test:	
	@echo "Testing..."
	@GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID} GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET} GOOGLE_REDIRECT_URL=${GOOGLE_REDIRECT_URL} SECRET_KEY=${SECRET_KEY} API_PORT=${API_PORT} DB_HOST=${DB_HOST} DB_PORT=${DB_PORT} DB_NAME=${DB_NAME} DB_USERNAME=${DB_USERNAME} DB_PASSWORD=${DB_PASSWORD} go test -p 1 -timeout 600s ./test/...
	

