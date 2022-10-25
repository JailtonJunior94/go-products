tests:
	go test -v --coverprofile test/coverage.out ./... 
	go tool cover -html=test/coverage.out

swagger:
	swag fmt
	swag init -g cmd/server/main.go