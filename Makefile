PROJECT_ROOT := $(shell sh -c "git rev-parse --path-format=absolute --git-dir | xargs dirname")

build:
	go build -o bin/converter cmd/converter/main.go

clean:
	rm -rvf ${PROJECT_ROOT}/bin
