package commands

import (
	"testing"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/stretchr/testify/assert"
)

func TestImageIsDangling(t *testing.T) {
	assert.True(t, (&Image{Image: image.Summary{
		Containers: 0,
		RepoTags:   nil,
	}}).IsDangling())
	assert.False(t, (&Image{Image: image.Summary{
		Containers: 0,
		RepoTags:   []string{"nginx:latest"},
	}}).IsDangling())
	assert.False(t, (&Image{Image: image.Summary{
		Containers: 1,
		RepoTags:   nil,
	}}).IsDangling())
}

func TestVolumeIsDangling(t *testing.T) {
	assert.True(t, (&Volume{Volume: &volume.Volume{
		UsageData: &volume.UsageData{RefCount: 0},
	}}).IsDangling())
	assert.False(t, (&Volume{Volume: &volume.Volume{
		UsageData: &volume.UsageData{RefCount: 2},
	}}).IsDangling())
	assert.False(t, (&Volume{Volume: &volume.Volume{}}).IsDangling())
}

func TestNetworkIsDangling(t *testing.T) {
	assert.True(t, (&Network{Network: network.Inspect{Containers: map[string]network.EndpointResource{}}}).IsDangling())
	assert.False(t, (&Network{Network: network.Inspect{
		Containers: map[string]network.EndpointResource{"abc": {}},
	}}).IsDangling())
}
