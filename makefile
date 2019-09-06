.PHONY: test serve

default: compose

deps:
	go get -v ./...

test: deps
	go test ./...

serve: test
	-rm -f ml main
	go install src/ml/main.go
	-@mv main ml 
	-@./ml || ml

local:
	export GOBIN="$$(pwd)" && make serve

compose:
	docker-compose up