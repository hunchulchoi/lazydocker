package presentation

import "github.com/jesseduffield/lazydocker/pkg/commands"

func GetNetworkDisplayStrings(network *commands.Network) []string {
	return displayStringsMutedIf([]string{network.Network.Driver, network.Name}, network.IsDangling())
}
