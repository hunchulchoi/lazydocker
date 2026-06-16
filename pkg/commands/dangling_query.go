package commands

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
)

func (c *DockerCommand) danglingImageIDs() map[string]struct{} {
	images, err := c.Client.ImageList(context.Background(), image.ListOptions{
		Filters: filters.NewArgs(filters.Arg("dangling", "true")),
	})
	if err != nil {
		return nil
	}

	ids := make(map[string]struct{}, len(images))
	for _, img := range images {
		ids[img.ID] = struct{}{}
	}

	return ids
}

func (c *DockerCommand) danglingVolumeNames() map[string]struct{} {
	result, err := c.Client.VolumeList(context.Background(), volume.ListOptions{
		Filters: filters.NewArgs(filters.Arg("dangling", "true")),
	})
	if err != nil {
		return nil
	}

	names := make(map[string]struct{}, len(result.Volumes))
	for _, vol := range result.Volumes {
		names[vol.Name] = struct{}{}
	}

	return names
}

func (c *DockerCommand) danglingNetworkNames() map[string]struct{} {
	networks, err := c.Client.NetworkList(context.Background(), network.ListOptions{
		Filters: filters.NewArgs(filters.Arg("dangling", "true")),
	})
	if err != nil {
		return nil
	}

	names := make(map[string]struct{}, len(networks))
	for _, nw := range networks {
		names[nw.Name] = struct{}{}
	}

	return names
}

func mapContains(set map[string]struct{}, key string) bool {
	if set == nil {
		return false
	}

	_, ok := set[key]
	return ok
}
