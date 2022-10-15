package main

import (
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func TestGatherItems(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	panicOnError(err)

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	panicOnError(err)
	defer reader.Close()
	io.Copy(io.Discard, reader)

	created, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "docker.io/library/alpine",
		Tty:   false,
	}, nil, nil, nil, "")
	panicOnError(err)

	err = cli.ContainerStart(ctx, created.ID, types.ContainerStartOptions{})
	panicOnError(err)

	cli.ContainerWait(ctx, created.ID, container.WaitConditionNotRunning)
	// panicOnError(err)

	items, err := gatherItems(cli)
	panicOnError(err)

	found := false
	for _, item := range items.Containers {
		if item.ID == created.ID {
			found = true
		}
	}
	err = cli.ContainerStop(ctx, created.ID, nil)
	panicOnError(err)

	err = cli.ContainerRemove(ctx, created.ID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
	panicOnError(err)

	if !found {
		t.Fatal("Created container not found")
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
