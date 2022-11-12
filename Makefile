generate-mock:
	go generate ./...

build: go build -o server.out -v ./cmd/server


run-coverage:
	go test Avito-Internship-Task/... -covermode=atomic -coverprofile=cover
	cat cover | fgrep -v "mocks" | fgrep -v "testing.go" | fgrep -v "docs"  | fgrep -v "configs" | fgrep -v "main.go" > cover2
	go tool cover -func=cover2

cover-html:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html

clean:
	rm -rf cover.html cover cover2 *.out
