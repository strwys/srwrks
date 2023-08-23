go-test:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

gen-mock:
	@mockery \
	--dir=internal/repository \
	--name=UserRepository \
	--filename=user_repository.go \
	--output=internal/mocks \
	--outpkg=mocks