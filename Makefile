.PHONY: default clean

default: server
all: default

export PATH := $(shell pwd)/tools/database-type/dist:$(PATH)

go_source = $(shell find src -type f -name "*.go")

tools: ./tools/database-type/database-type.go
	@echo -e "   GO      build ./tools/database-type/database-type.go"
	@go build -o ./tools/database-type/dist/database-type ./tools/database-type/database-type.go

server: tools $(go_source)
	@echo -e "   GO      generate ./..."
	@go generate ./...
	@echo -e "   GO      build ./src/main.go"
	@go build -o dist/server ./src/main.go

clean:
	@echo -e "   RM      dist/server"
	@rm -f ./dist/server
	@echo -e "   RM      tools/database-type/dist"
	@rm -Rf tools/database-type/dist
