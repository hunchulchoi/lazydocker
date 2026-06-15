package commands

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
)

// Volume : A docker Volume
type Volume struct {
	Name          string
	Volume        *volume.Volume
	Client        *client.Client
	OSCommand     *OSCommand
	Log           *logrus.Entry
	DockerCommand LimitedDockerCommand
}

// RefreshVolumes gets the volumes and stores them
func (c *DockerCommand) RefreshVolumes() ([]*Volume, error) {
	result, err := c.Client.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	usageByName := c.volumeUsageByName()
	volumes := result.Volumes

	ownVolumes := make([]*Volume, len(volumes))

	for i, vol := range volumes {
		if usageByName != nil {
			if usageData, ok := usageByName[vol.Name]; ok {
				vol.UsageData = usageData
			}
		}

		ownVolumes[i] = &Volume{
			Name:          vol.Name,
			Volume:        vol,
			Client:        c.Client,
			OSCommand:     c.OSCommand,
			Log:           c.Log,
			DockerCommand: c,
		}
	}

	return ownVolumes, nil
}

func (c *DockerCommand) volumeUsageByName() map[string]*volume.UsageData {
	usage, err := c.Client.DiskUsage(context.Background(), types.DiskUsageOptions{
		Types: []types.DiskUsageObject{types.VolumeObject},
	})
	if err != nil {
		return nil
	}

	usageByName := make(map[string]*volume.UsageData, len(usage.Volumes))
	for _, vol := range usage.Volumes {
		if vol != nil && vol.UsageData != nil {
			usageByName[vol.Name] = vol.UsageData
		}
	}

	return usageByName
}

// VolumeSizesByName returns volume disk usage keyed by volume name.
func (c *DockerCommand) VolumeSizesByName() (map[string]int64, error) {
	usageByName := c.volumeUsageByName()
	if usageByName == nil {
		return nil, errors.New("failed to get volume disk usage")
	}

	sizes := make(map[string]int64, len(usageByName))
	for name, usageData := range usageByName {
		sizes[name] = usageData.Size
	}

	return sizes, nil
}

// PruneVolumes prunes volumes
func (c *DockerCommand) PruneVolumes() error {
	_, err := c.Client.VolumesPrune(context.Background(), filters.Args{})
	return err
}

// Remove removes the volume
func (v *Volume) Remove(force bool) error {
	return v.Client.VolumeRemove(context.Background(), v.Name, force)
}
