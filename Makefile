APP_NAME?=my-app

# сборка отдельного приложения
clean:
	rm -f ${APP_NAME}

build: clean
	go build -mod=mod -o ${APP_NAME} ./cmd/service/service.go

run: build
	./${APP_NAME}

deploy: create-migrator
	docker-compose up --build
	sleep 4
	docker run --network host migrator  \
	-path=/migrations/ \
	-database "postgresql://postgres:postgres@localhost:7557/vk?sslmode=disable" up
	
.PHONY: test
test:
	go test -v -race -count=1 ./...

.PHONY: gen_docs
gen_docs:
	/home/andreysm/go/bin/swag init -g ./cmd/service/service.go

.PHONY: create-migrator
create-migrator:
	docker build -t migrator ./db

.PHONY: shutdown
shutdown:
	docker-compose down -v