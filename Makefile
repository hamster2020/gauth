.PHONY: check check-docker createdb dropdb install mocks

check:
	@go fmt ./...
	@go vet ./...
	@go test -cover -race ./... 

check-docker:
	@go fmt ./...
	@go vet ./...
	@go test -cover -race `go list ./... | grep -v postgres | grep -v password-validator`

createdb:
	@createdb gauth || true

dropdb:
	@dropdb gauth || true

install:
	@go install ./...

mocks:
	rm mocks/* || true
	mockit 'ds *mockDatastore' Datastore mocks > mocks/datastore_mock_gen.go
	mockit 'token *MockToken' Token mocks > mocks/token_mock_gen.go
	mockit 'logic *MockLogic' Logic mocks > mocks/logic_mock_gen.go
	mockit 'validator *mockValidator' Validator mocks > mocks/validator_mock_gen.go
	goimports -w ./mocks/*
	gofmt -w ./mocks/*