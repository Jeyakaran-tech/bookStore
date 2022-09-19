export DB_USER=user 
export DB_PASS=testbooks
export INSTANCE_CONNECTION_NAME=bookstore-362511:australia-southeast2:bookstore
export DB_NAME=books

compile:
	go build -v -o build/bookstore ./cmd/app

run:
	go run cmd/app/main.go

test:
	go test -v -race -timeout 1000s -covermode=atomic -coverpkg=./integrationtests -coverprofile=unit_test.raw.out ./...

integration_test:
	go test -v -race -timeout 1000s -covermode=atomic -coverpkg=./integrationtests -coverprofile=unit_test.raw.out ./...
