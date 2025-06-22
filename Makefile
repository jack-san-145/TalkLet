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

