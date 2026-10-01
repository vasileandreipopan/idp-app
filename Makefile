# says these are not real targets, just names for commands to run
.PHONY: fmt vet run

# format the code in the current directory and all subdirectories
fmt:
	gofmt -w .

# run go vet on the current directory and all subdirectories
vet:
	go vet ./...

# run the fmt and vet targets, then run the main.go file
run: fmt vet
	go run .