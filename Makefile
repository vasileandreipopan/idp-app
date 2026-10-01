# says these are not real targets, just names for commands to run
.PHONY: fmt vet run docker-build docker-run

# sets this only if is not already set in the env
VERSION   ?= $(shell git describe --always --dirty)
IMAGE     ?= idp-app
TAG       ?= local
HOST_PORT ?= 8888

# format the code in the current directory and all subdirectories
fmt:
	gofmt -w .

# run go vet on the current directory and all subdirectories
vet:
	go vet ./...

# run the fmt and vet targets, then run the main.go file
run: fmt vet
	go run .

# build a docker image with the specified version and tag
docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(TAG) .

# run a docker container with the specified image and tag, mapping the host port to the container port
docker-run:
	docker run --rm -p $(HOST_PORT):8888 $(IMAGE):$(TAG)