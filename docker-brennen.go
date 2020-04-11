package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/urfave/cli/v2"
)

type item struct {
	ID          string
	Description string
}

type items struct {
	Containers []item
	Images     []item
	Networks   []item
	Volumes    []item
}

func main() {
	app := &cli.App{
		Name:            "docker-brennen",
		Usage:           "cleanup unused Docker resources",
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			return run()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	docker, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	items := items{}

	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: singleArg("status", "exited")})
	if err != nil {
		return err
	}
	for _, container := range containers {
		items.Containers = append(items.Containers, item{
			ID:          container.ID,
			Description: strings.Join(container.Names, ", "),
		})
	}

	images, err := docker.ImageList(context.Background(), types.ImageListOptions{Filters: singleArg("dangling", "true")})
	if err != nil {
		return err
	}
	for _, image := range images {
		items.Images = append(items.Images, item{
			ID:          image.ID[7:],
			Description: strings.Join(image.RepoTags, ", "),
		})
	}

	networks, err := docker.NetworkList(context.Background(), types.NetworkListOptions{Filters: singleArg("driver", "bridge")})
	if err != nil {
		return err
	}
	for _, network := range networks {
		if network.Name != "bridge" && len(network.Containers) == 0 {
			items.Networks = append(items.Networks, item{
				ID:          network.ID,
				Description: network.Name + "/" + network.Driver,
			})
		}
	}

	volumeList, err := docker.VolumeList(context.Background(), singleArg("dangling", "true"))
	if err != nil {
		return err
	}
	for _, volume := range volumeList.Volumes {
		items.Volumes = append(items.Volumes, item{
			ID:          volume.Name,
			Description: volume.Mountpoint,
		})
	}

	containerCount := len(items.Containers)
	imageCount := len(items.Images)
	networkCount := len(items.Networks)
	volumeCount := len(items.Volumes)

	if containerCount+imageCount+networkCount+volumeCount == 0 {
		fmt.Println("Nothing to remove")
	} else {
		fmt.Println("TYPE       ID            DESCRIPTION")
		for _, item := range items.Containers {
			fmt.Printf("%s  %s  %s\n", "container", item.ID[:12], item.Description)
		}
		for _, item := range items.Images {
			fmt.Printf("%s  %s  %s\n", "image    ", item.ID[:12], item.Description)
		}
		for _, item := range items.Networks {
			fmt.Printf("%s  %s  %s\n", "network  ", item.ID[:12], item.Description)
		}
		for _, item := range items.Volumes {
			fmt.Printf("%s  %s  %s\n", "volume   ", item.ID[:12], item.Description)
		}

		fmt.Printf("Are you sure you want to remove %d containers, %d images, %d networks and %d volumes? [y/n]\n",
			containerCount,
			imageCount,
			networkCount,
			volumeCount,
		)

		var response string
		_, err = fmt.Scanln(&response)
		if err != nil {
			return err
		}
		if response == "y" {
			for _, item := range items.Containers {
				removeContainer(docker, item)
			}
			for _, item := range items.Images {
				removeImage(docker, item)
			}
			for _, item := range items.Networks {
				removeNetwork(docker, item)
			}
			for _, item := range items.Volumes {
				removeVolume(docker, item)
			}
		} else {
			fmt.Println("Nothing has been removed")
		}
	}
	return nil
}

func singleArg(name string, value string) filters.Args {
	args := filters.NewArgs()
	args.Add(name, value)
	return args
}

func removeContainer(docker *client.Client, item item) {
	err := docker.ContainerRemove(context.Background(), item.ID, types.ContainerRemoveOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s removed\n", "Container", item.ID[:12])
}

func removeImage(docker *client.Client, item item) {
	_, err := docker.ImageRemove(context.Background(), item.ID, types.ImageRemoveOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s removed\n", "Image", item.ID[:12])
}

func removeNetwork(docker *client.Client, item item) {
	err := docker.NetworkRemove(context.Background(), item.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s removed\n", "Network", item.ID[:12])
}

func removeVolume(docker *client.Client, item item) {
	err := docker.VolumeRemove(context.Background(), item.ID, false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s removed\n", "Volume", item.ID[:12])
}
