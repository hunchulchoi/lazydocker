package presentation

import "github.com/jesseduffield/lazydocker/pkg/commands"

func GetVolumeDisplayStrings(volume *commands.Volume) []string {
	return displayStringsMutedIf([]string{volume.Volume.Driver, volume.Name}, volume.IsDangling())
}
