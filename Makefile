all: dev

dev:
	go run main.go&
	npm run dev&

prod:
	npm run build && go build

kill:
	@lsof -t -i:8080 && kill -9 $$(lsof -t -i:8080)
	@lsof -t -i:3000 && kill -9 $$(lsof -t -i:3000)