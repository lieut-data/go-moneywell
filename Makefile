.PHONY: test coverage

test:
	go test ./... -coverprofile=coverage.out 

coverage: test
	go tool cover -html=coverage.out
	
