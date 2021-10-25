.PHONY: check install mocks

check:
	@go fmt ./...
	@go vet ./...
	@go test -cover -race ./... 

install:
	@go install ./...

mocks:
	rm mocks/* || true
	mockit 'ds *mockDatastore' Datastore mocks > mocks/datastore_mock_gen.go
	mockit 'logic *MockLogic' Logic mocks > mocks/logic_mock_gen.go
	mockit 'validator *mockValidator' Validator mocks > mocks/validator_mock_gen.go
	goimports -w ./mocks/*
	gofmt -w ./mocks/*