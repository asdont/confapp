build:
	swag init -g cmd/app/main.go
	go build -o confapp cmd/app/main.go

run:
	swag init -g cmd/app/main.go
	go run cmd/app/main.go

d.build:
	swag init -g cmd/app/main.go
	docker build . -t confapp:v0.1

d.run.localdb:
	docker run --name confapp --network=host confapp:v0.1

d.start:
	docker start confapp

d.c.build:
	swag init -g cmd/app/main.go
	docker-compose build

d.c.run:
	docker-compose up