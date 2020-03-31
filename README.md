# docker-brennen

Cleanup unused Docker resources:

- Exited containers.
- Dangling images.
- Non-bridge networks.
- Dangling volumes.

## Usage

### Running

- List all Docker resources to remove and ask for confirmation (type `y`) before removing them:

```shell
docker-brennen
```

## Development

### Building

- This is a simple standard [Go](https://golang.org/) project.
- Go version 1.14 is required.
- Build using `go build`.
- Run with `go run .`.
