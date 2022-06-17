run: migrate
	go run github.com/drTragger/RileyBot/cmd/bot
migrate:
	migrate -path migrations -database "mysql://$(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_DATABASE)" up
rollback:
	migrate -path migrations -database "mysql://$(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_DATABASE)" down
