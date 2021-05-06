.PHONY: start build


APP_NAME		= echo-admin
BUILD_ROOT		= build

all: start

fmt:
	@find . -name "*.go" -type f -not -path "./vendor/*"|xargs gofmt -s -w

build:
	@go build -ldflags "-w -s" -o $(BUILD_ROOT)/$(APP_NAME)

start:
	@go run ./main.go runserver --config=./config/config.yaml --casbin_model=./config/casbin_model.conf

migrate:
	@go run ./main.go migrate --config=./config/config.yaml

setup:
	@go run ./main.go setup --config=./config/config.yaml --menu=./config/menu.yaml

swagger:
	@swag init --parseDependency --parseInternal -g api/routes/swagger_route.go

clean:
	@rm -rf $(BUILD_ROOT)
