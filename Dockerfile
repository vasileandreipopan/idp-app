# FIRST STAGE: Build
# start from the official go image, which includes the Go compiler and standard library
FROM golang:1.27 AS build

WORKDIR /src

# this copies these two files from the laptop into /src in the image
COPY go.mod main.go ./

# build argument to specify the version of the application, defaulting to "dev" if not provided
ARG VERSION=dev

# executes a command while building the image
# CGO_ENABLED=0 builds a static binary 
# -trimpath removes all file system paths from the compiled executable
# -ldflags="" is used to pass options to the linker, in this case:
# -s -w strips the symbol table and debug information from the binary
# -X main.version=${VERSION} sets the value of the version variable in the main package
# -o /out/idp-app specifies the output file for the compiled binary
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=${VERSION}" -o /out/idp-app .


# FINAL STAGE: Run
# gcr.io/distroless/static is a distroless image for static binaries
# nonroot is a variant of the distroless image that runs as a non-root user
FROM gcr.io/distroless/static:nonroot AS final

# copy the compiled binary from the build stage into the final image
COPY --from=build /out/idp-app /idp-app

# run the app as a non-root user
USER nonroot:nonroot
EXPOSE 8888

ENTRYPOINT ["/idp-app"]