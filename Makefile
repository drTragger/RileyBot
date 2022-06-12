run: migrate
	go run github.com/drTragger/RileyBot/cmd/bot
migrate:
	migrate -path migrations -database "mysql://misha:1m2i3s4h5a@tcp(localhost:3306)/riley" up
rollback:
	migrate -path migrations -database "mysql://misha:1m2i3s4h5a@tcp(localhost:3306)/riley" down
