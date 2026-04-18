include .env
export

migrate_up:
	migrate -path ./internals/storage/postgres/migrations -database "$(POSTGRES_DATABASE_CONNECTION)" up

migrate_down:
	migrate -path ./internals/storage/postgres/migrations -database "$(POSTGRES_DATABASE_CONNECTION)" down

build:
	cd cmd/server && go build -o talklet && ./talklet
	
run:
	cd cmd/server && ./talklet

start:
	docker start c8d7cf87cf78 && docker start c97328d35a6d

stop:
	docker stop c8d7cf87cf78 && docker stop c97328d35a6d