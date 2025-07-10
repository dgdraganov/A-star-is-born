
build:
	go build -o bin/astar ./cmd/astar/main.go

run: build
	./bin/astar