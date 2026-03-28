# gosmos

This Go package wraps Typst, a modern typesetting system, in a Go package. The only public interface is the `typst` program. which this package executes.

It does not expect the binary to be present on the executing host. Instead, this package provides a Go tool that fetches the Typst binary from the official release page using a `go-generate` statement.

The Go program that uses this package is assumed to be packaged using [ko](https://ko.build/features/static-assets/), which treats the Typst binary as static asset. This is a compromise betweeen the need to have the binary present on the host system and a separate build step for different platforms, or even embedding it into the Go binary.

# Usage

1. Add this package to your Go project by running:

    ```sh
    go get github.com/suhlig/gosmos
    ```

1. Get the Go tool for fetching the typst binary at build time:

    ```command
    go get -tool github.com/suhlig/gosmos/cmd/fetch-typst
    ```

1. Add the code that fetches the Typst binary for inclusion as a static asset:

    ```go
    //go:generate go tool github.com/suhlig/gosmos/cmd/fetch-typst --path kodata/typst
    ```

1. Build and run your container on a Mac:

    ```command
    docker run -p 8080:8080 $(ko build --local --platform=linux/arm64 ./cmd/server)
    ```

1. Push the container to a remote host using [docker pussh](https://github.com/psviderski/unregistry):

    ```command
    docker pussh --platform linux/amd64 $(ko build --local --platform=linux/arm64 ./cmd/server) alice@my.docker.host.example.com
    ```

1. Take note of the image name and run the container on the remote host using:

    ```command
    docker run -p 8090:8080 --rm -it ko.local/server-bfaa5a9da622cb30b1dc28c3904f0eee:cabc6c5ae22a54fba9793d084c7bbd8957cf45ba81191a76b0aae5297bd27bc1
    ```

# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
