package presentation

import "github.com/hunchulchoi/lazydocker/pkg/commands"

func GetVolumeDisplayStrings(volume *commands.Volume) []string {
	return displayStringsMutedIf([]string{volume.Volume.Driver, volume.Name}, volume.IsDangling())
}
