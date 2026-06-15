package commands

import (
	"sort"

	"github.com/docker/docker/api/types/container"
	dockermount "github.com/docker/docker/api/types/mount"
	"github.com/jesseduffield/lazydocker/pkg/i18n"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

// MountSizes provides optional disk usage data for container mounts.
type MountSizes struct {
	VolumeSizes map[string]int64
	OSCommand   *OSCommand
}

func (c *Container) RenderVolumes(tr *i18n.TranslationSet, sizes MountSizes) (string, error) {
	if !c.DetailsLoaded() {
		return tr.WaitingForContainerInfo, nil
	}

	if len(c.Details.Mounts) == 0 {
		return tr.NothingToDisplay, nil
	}

	rows := [][]string{
		{
			tr.VolumesTypeHeader,
			tr.VolumesNameHeader,
			tr.VolumesSourceHeader,
			tr.VolumesDestinationHeader,
			tr.VolumesSizeHeader,
			tr.VolumesDriverHeader,
			tr.VolumesModeHeader,
			tr.VolumesAccessHeader,
		},
	}

	for _, row := range mountRows(c.Details.Mounts, sizes) {
		rows = append(rows, row)
	}

	return utils.RenderTable(rows)
}

func mountRows(mounts []container.MountPoint, sizes MountSizes) [][]string {
	sorted := append([]container.MountPoint(nil), mounts...)
	sort.Slice(sorted, func(i, j int) bool {
		left := string(sorted[i].Type) + sorted[i].Destination
		right := string(sorted[j].Type) + sorted[j].Destination
		return left < right
	})

	rows := make([][]string, 0, len(sorted))
	for _, mount := range sorted {
		rows = append(rows, []string{
			string(mount.Type),
			displayOrDash(mount.Name),
			displayOrDash(mount.Source),
			displayOrDash(mount.Destination),
			formatMountSize(mount, sizes),
			displayOrDash(mount.Driver),
			displayOrDash(mount.Mode),
			mountAccess(mount.RW),
		})
	}

	return rows
}

func formatMountSize(m container.MountPoint, sizes MountSizes) string {
	switch m.Type {
	case dockermount.TypeVolume:
		if sizes.VolumeSizes == nil {
			return "-"
		}
		size, ok := sizes.VolumeSizes[m.Name]
		if !ok || size < 0 {
			return "-"
		}
		return utils.FormatBinaryBytes(int(size))
	case dockermount.TypeBind:
		if sizes.OSCommand == nil || m.Source == "" {
			return "-"
		}
		size, err := sizes.OSCommand.DirSizeBytes(m.Source)
		if err != nil {
			return "-"
		}
		return utils.FormatBinaryBytes(int(size))
	default:
		return "-"
	}
}

func mountAccess(readWrite bool) string {
	if readWrite {
		return "rw"
	}
	return "ro"
}

func displayOrDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
