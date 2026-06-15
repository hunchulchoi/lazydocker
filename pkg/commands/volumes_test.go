package commands

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/stretchr/testify/assert"
)

func TestRenderVolumes(t *testing.T) {
	tr := testTranslationSet(t)

	c := &Container{
		Details: container.InspectResponse{
			ContainerJSONBase: &container.ContainerJSONBase{},
			Mounts: []container.MountPoint{
				{
					Type:        mount.TypeVolume,
					Name:        "data",
					Source:      "/var/lib/docker/volumes/data/_data",
					Destination: "/var/lib/mysql",
					Driver:      "local",
					Mode:        "z",
					RW:          true,
				},
				{
					Type:        mount.TypeBind,
					Source:      "/host/config",
					Destination: "/etc/app",
					Mode:        "ro",
					RW:          false,
				},
			},
		},
	}

	sizes := MountSizes{
		VolumeSizes: map[string]int64{
			"data": 2048,
		},
	}

	output, err := c.RenderVolumes(tr, sizes)
	assert.NoError(t, err)
	assert.Contains(t, output, "bind")
	assert.Contains(t, output, "/host/config")
	assert.Contains(t, output, "volume")
	assert.Contains(t, output, "data")
	assert.Contains(t, output, "2.00kiB")
	assert.Contains(t, output, "rw")
	assert.Contains(t, output, "ro")
}

func TestRenderVolumesNothingToDisplay(t *testing.T) {
	tr := testTranslationSet(t)
	c := &Container{
		Details: container.InspectResponse{
			ContainerJSONBase: &container.ContainerJSONBase{},
			Mounts:            nil,
		},
	}

	output, err := c.RenderVolumes(tr, MountSizes{})
	assert.NoError(t, err)
	assert.Equal(t, tr.NothingToDisplay, output)
}

func TestFormatMountSizeVolumeUnavailable(t *testing.T) {
	size := formatMountSize(container.MountPoint{
		Type: mount.TypeVolume,
		Name: "missing",
	}, MountSizes{
		VolumeSizes: map[string]int64{
			"other": 100,
		},
	})
	assert.Equal(t, "-", size)
}

func TestFormatMountSizeVolumeNegative(t *testing.T) {
	size := formatMountSize(container.MountPoint{
		Type: mount.TypeVolume,
		Name: "remote",
	}, MountSizes{
		VolumeSizes: map[string]int64{
			"remote": -1,
		},
	})
	assert.Equal(t, "-", size)
}
