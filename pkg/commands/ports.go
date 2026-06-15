package commands

import (
	"sort"
	"strings"

	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/jesseduffield/lazydocker/pkg/i18n"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

func (c *Container) RenderPorts(tr *i18n.TranslationSet) (string, error) {
	if !c.DetailsLoaded() {
		return tr.WaitingForContainerInfo, nil
	}

	networks := containerNetworkNames(c.Details.NetworkSettings.Networks)
	networksLabel := strings.Join(networks, ", ")
	if networksLabel == "" {
		networksLabel = "-"
	}

	rows := [][]string{
		{
			tr.PortsInternalHeader,
			tr.PortsExternalHeader,
			tr.PortsHostHeader,
			tr.PortsNetworkHeader,
		},
	}

	portRows := publishedPortRows(c.Details.NetworkSettings.Ports, networksLabel)
	if len(portRows) == 0 {
		portRows = exposedPortRows(c.Details.Config.ExposedPorts, networksLabel)
	}
	rows = append(rows, portRows...)

	if len(portRows) == 0 {
		return tr.NothingToDisplay, nil
	}

	table, err := utils.RenderTable(rows)
	if err != nil {
		return "", err
	}

	networkSection := renderNetworkSection(tr, c.Details.NetworkSettings.Networks)
	if networkSection == "" {
		return table, nil
	}

	return table + "\n\n" + networkSection, nil
}

func publishedPortRows(ports nat.PortMap, networksLabel string) [][]string {
	if len(ports) == 0 {
		return nil
	}

	portKeys := make([]nat.Port, 0, len(ports))
	for port := range ports {
		portKeys = append(portKeys, port)
	}
	sort.Slice(portKeys, func(i, j int) bool {
		return string(portKeys[i]) < string(portKeys[j])
	})

	rows := make([][]string, 0, len(portKeys))
	for _, port := range portKeys {
		bindings := ports[port]
		if len(bindings) == 0 {
			rows = append(rows, []string{string(port), "-", "-", networksLabel})
			continue
		}

		for _, binding := range bindings {
			external := binding.HostPort
			if external == "" {
				external = "-"
			}
			hostIP := binding.HostIP
			if hostIP == "" {
				hostIP = "0.0.0.0"
			}
			rows = append(rows, []string{string(port), external, hostIP, networksLabel})
		}
	}

	return rows
}

func exposedPortRows(exposed nat.PortSet, networksLabel string) [][]string {
	if len(exposed) == 0 {
		return nil
	}

	portKeys := make([]nat.Port, 0, len(exposed))
	for port := range exposed {
		portKeys = append(portKeys, port)
	}
	sort.Slice(portKeys, func(i, j int) bool {
		return string(portKeys[i]) < string(portKeys[j])
	})

	rows := make([][]string, 0, len(portKeys))
	for _, port := range portKeys {
		rows = append(rows, []string{string(port), "-", "-", networksLabel})
	}

	return rows
}

func containerNetworkNames(networks map[string]*network.EndpointSettings) []string {
	if len(networks) == 0 {
		return nil
	}

	names := make([]string, 0, len(networks))
	for name := range networks {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func renderNetworkSection(tr *i18n.TranslationSet, networks map[string]*network.EndpointSettings) string {
	if len(networks) == 0 {
		return ""
	}

	names := containerNetworkNames(networks)
	rows := [][]string{
		{
			tr.PortsNetworkHeader,
			tr.PortsContainerIPHeader,
			tr.PortsGatewayHeader,
		},
	}

	for _, name := range names {
		endpoint := networks[name]
		ip := "-"
		gateway := "-"
		if endpoint != nil {
			if endpoint.IPAddress != "" {
				ip = endpoint.IPAddress
			}
			if endpoint.Gateway != "" {
				gateway = endpoint.Gateway
			}
		}
		rows = append(rows, []string{name, ip, gateway})
	}

	table, err := utils.RenderTable(rows)
	if err != nil {
		return ""
	}

	return tr.PortsNetworksSection + "\n" + table
}
