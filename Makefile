APP_NAME?=my-app

# сборка отдельного приложения
clean:
	rm -f ${APP_NAME}

build: clean
	go build -o ${APP_NAME} ./cmd/service/service.go

run: build
	./${APP_NAME}