export MYSQL_USER=user 
export MYSQL_PASSWORD=testbooks
export MYSQL_INSTANCE=bookstore-362511:australia-southeast2:bookstore
export MYSQL_DATABASE=books
export MYSQL_PORT=5432
export MYSQL_HOST=127.0.0.1

compile:
	go build -v -o build/bookstore ./cmd/app

run:
	go run cmd/app/main.go

test:
	go test -v -race -timeout 1000s -covermode=atomic -coverpkg=./cloudsql -coverprofile=unit_test.raw.out ./...
