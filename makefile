BINARY=engine

test: 
	go test -v -cover -covermode=atomic ./...

engine:
	go build -o bin/${BINARY} main.go

unittest:
	go test -short  ./...

clean:
	if [ -f bin/${BINARY} ] ; then rm bin/${BINARY} ; fi

docker:
	docker build -t consul-demo-app .

run:
	docker-compose up --build -d

stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: clean install unittest build docker run stop vendor lint-prepare lint