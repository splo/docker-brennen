package main

import (
	"bufio"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func TestGatherItems(t *testing.T) {
	tc, cancel := BuildTestingContext(t)
	defer cancel()
	tc.PullImage("docker.io/library/alpine")
	id := tc.CreateContainer("docker.io/library/alpine")
	tc.StartContainer(id)
	tc.StopContainer(id)
	tc.WaitContainerStopped(id)
	items, err := gatherItems(tc.cli)
	if err != nil {
		t.Fatal(err)
	}
	found := ItemsContainId(items.Containers, id)
	tc.RemoveContainer(id)
	if found {
		t.Log("Created container found")
	} else {
		t.Fatal("Created container not found")
	}
}

func ItemsContainId(items []item, ID string) bool {
	for _, item := range items {
		if item.ID == ID {
			return true
		}
	}
	return false
}

func BuildTestingContext(t *testing.T) (TestingContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatal(err)
	}
	return TestingContext{t: t, ctx: ctx, cli: cli}, cancel
}

type TestingContext struct {
	t   *testing.T
	ctx context.Context
	cli *client.Client
}

func (tc *TestingContext) PullImage(image string) {
	reader, err := tc.cli.ImagePull(tc.ctx, image, types.ImagePullOptions{})
	if err != nil {
		tc.t.Fatal(err)
	}
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		var jsonLine map[string]any
		err = json.Unmarshal([]byte(line), &jsonLine)
		if err != nil {
			tc.t.Log(line)
		} else {
			if jsonLine["progress"] != nil {
				tc.t.Log(jsonLine["status"], jsonLine["progress"])
			} else {
				tc.t.Log(jsonLine["status"])
			}
		}
	}
}

func (tc *TestingContext) CreateContainer(image string) string {
	created, err := tc.cli.ContainerCreate(tc.ctx, &container.Config{
		Image: image,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		tc.t.Fatal(err)
	}
	tc.t.Log("Container", created.ID, "created")
	return created.ID
}

func (tc *TestingContext) StartContainer(ID string) {
	err := tc.cli.ContainerStart(tc.ctx, ID, container.StartOptions{})
	if err != nil {
		tc.t.Fatal(err)
	}
	tc.t.Log("Container", ID, "started")
}

func (tc *TestingContext) StopContainer(ID string) {
	err := tc.cli.ContainerStop(tc.ctx, ID, container.StopOptions{})
	if err != nil {
		tc.t.Fatal(err)
	}
	tc.t.Log("Container", ID, "stop requested")
}

func (tc *TestingContext) WaitContainerStopped(ID string) {
	statusCh, errCh := tc.cli.ContainerWait(tc.ctx, ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			tc.t.Fatal(err)
		}
	case <-statusCh:
	}
	tc.t.Log("Container", ID, "stopped")
}

func (tc *TestingContext) RemoveContainer(ID string) {
	err := tc.cli.ContainerRemove(tc.ctx, ID, types.ContainerRemoveOptions{RemoveVolumes: true, Force: true})
	if err != nil {
		tc.t.Fatal(err)
	}
	tc.t.Log("Container removed")
}
