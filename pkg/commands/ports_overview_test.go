package commands

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
)

func TestCollectAllPortsConflict(t *testing.T) {
	port, err := nat.NewPort("tcp", "80")
	assert.NoError(t, err)

	containers := []*Container{
		{
			Name: "web-a",
			Container: container.Summary{
				State: "running",
			},
			Details: container.InspectResponse{
				ContainerJSONBase: &container.ContainerJSONBase{},
				Config:            &container.Config{},
				NetworkSettings: &container.NetworkSettings{
					NetworkSettingsBase: container.NetworkSettingsBase{
						Ports: nat.PortMap{
							port: []nat.PortBinding{
								{HostIP: "0.0.0.0", HostPort: "8080"},
							},
						},
					},
					Networks: map[string]*network.EndpointSettings{
						"bridge": {IPAddress: "172.17.0.2"},
					},
				},
			},
		},
		{
			Name: "web-b",
			Container: container.Summary{
				State: "running",
			},
			Details: container.InspectResponse{
				ContainerJSONBase: &container.ContainerJSONBase{},
				Config:            &container.Config{},
				NetworkSettings: &container.NetworkSettings{
					NetworkSettingsBase: container.NetworkSettingsBase{
						Ports: nat.PortMap{
							port: []nat.PortBinding{
								{HostIP: "0.0.0.0", HostPort: "8080"},
							},
						},
					},
					Networks: map[string]*network.EndpointSettings{
						"bridge": {IPAddress: "172.17.0.3"},
					},
				},
			},
		},
	}

	rows := CollectAllPorts(containers)
	assert.Len(t, rows, 2)
	assert.True(t, rows[0].Conflict)
	assert.True(t, rows[1].Conflict)
}

func TestSortPortOverviewRowsByExternal(t *testing.T) {
	rows := []PortOverviewRow{
		{Container: "b", External: "9000"},
		{Container: "a", External: "8080"},
	}

	SortPortOverviewRows(rows, PortSortExternal, true)
	assert.Equal(t, "8080", rows[0].External)
	assert.Equal(t, "9000", rows[1].External)
}

func TestPortOverviewHeaders(t *testing.T) {
	headers := PortOverviewHeaders([]string{"A", "B", "C", "D", "E", "F", "G"}, PortSortExternal, true)
	assert.Equal(t, "E ↑", headers[4])
}
