.PHONY: clean backend-build

GOFILES := $(wildcard app/backend/*.go)

docker:
	@docker build -t cms .

build: backend-build
	@mkdir -p out

backend-build: $(GOFILES)
	@$(shell cd app/backend && go build -o ../../out/cms ./main.go)

