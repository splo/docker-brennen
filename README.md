# docker-brennen

Cleanup unused Docker resources:

- Exited containers.
- Dangling images.
- Non-bridge networks.
- Dangling volumes.

## Usage

### Install

- On Macos with Homebrew:

```shell
brew install splo/tap/docker-brennen
```

- On other OSes by download:

  - Download an archive from <https://github.com/splo/docker-brennen/releases/latest>.
  - Unzip it.
  - Move to somewhere in your `$PATH`.

### Running

- Get help:

```shell
$ docker-brennen --help
NAME:
   docker-brennen - cleanup unused Docker resources

USAGE:
   docker-brennen [global options] [arguments...]

GLOBAL OPTIONS:
   --force, -f  remove resources without confirmation prompt (default: false)
   --help, -h   show help (default: false)
```

- List all Docker resources to remove and ask for confirmation (type `y`) before removing them:

```shell
$ docker-brennen
TYPE       ID            DESCRIPTION
container  a8173677c544  /festive_margulis
container  c4dba3072af2  /vigilant_poincare
container  cb77822a22fe  /friendly_mclaren
container  2c90f4bada6e  /nervous_easley
container  813aef4d04ed  /wonderful_merkle
container  3444377e9cab  /naughty_mcnulty
container  9d7462c15ab4  /unruffled_jackson
image      965978555d82  <none>:<none>
network    9804c17fa710  foo/bridge
volume     94bb96074976  /var/lib/docker/volumes/94bb9607497623326ce29a9aa1fdcc3054c7d6d248b6e7c60326554817a6e184/_data
Are you sure you want to remove 7 containers, 1 images, 1 networks and 1 volumes? [y/n]
y
Container a8173677c544 removed
Container c4dba3072af2 removed
Container cb77822a22fe removed
Container 2c90f4bada6e removed
Container 813aef4d04ed removed
Container 3444377e9cab removed
Container 9d7462c15ab4 removed
Image 965978555d82 removed
Network 9804c17fa710 removed
Volume 94bb96074976 removed
```

## Development

### Building

- This is a simple standard [Go](https://golang.org/) project.
- Go version 1.14 is required.
- Build using `go build`.
- Run with `go run .`.
