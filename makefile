kill:
	@lsof -ti :8000 | xargs kill -9

dev:
	air -c .air.toml

run:
	go run main.go