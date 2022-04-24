.PHONY: build

build: parse.go
	go build

parse.go: parse.y
	./golemon parse.y
