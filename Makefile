run:
	go run .

test:
	go test -v ./...

test-cover:
	go test -v -race -cover ./...