test-unit:
	go test -tags unit -v ./...

cover:
	go test -tags unit -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out