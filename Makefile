run:
	docker-compose down -v
	service mysql stop
	docker build --no-cache --network host -f ./docker/Dockerfile . --tag app
	docker-compose up --build

build:	generate-api
	go build -o server.out -v ./cmd/server/main.go

make-mocks:
	go generate ./...

generate-api:
	go install github.com/swaggo/swag/cmd/swag@v1.6.5
	swag init -g ./cmd/server/main.go -o docs

tests:	make-mocks
	go test ./...
	
clean:
	rm -rf *.out *.exe *.html
