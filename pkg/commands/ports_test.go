package commands

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/jesseduffield/lazydocker/pkg/i18n"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func testTranslationSet(t *testing.T) *i18n.TranslationSet {
	tr, err := i18n.NewTranslationSetFromConfig(logrus.NewEntry(logrus.New()), "en")
	assert.NoError(t, err)
	return tr
}

func TestRenderPortsPublished(t *testing.T) {
	tr := testTranslationSet(t)
	port, err := nat.NewPort("tcp", "80")
	assert.NoError(t, err)

	c := &Container{
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
					"bridge": {IPAddress: "172.17.0.2", Gateway: "172.17.0.1"},
					"app":    {IPAddress: "172.18.0.3", Gateway: "172.18.0.1"},
				},
			},
		},
	}

	output, err := c.RenderPorts(tr)
	assert.NoError(t, err)
	assert.Contains(t, output, "80/tcp")
	assert.Contains(t, output, "8080")
	assert.Contains(t, output, "app, bridge")
	assert.Contains(t, output, "172.18.0.3")
}

func TestRenderPortsExposedOnly(t *testing.T) {
	tr := testTranslationSet(t)
	port, err := nat.NewPort("tcp", "443")
	assert.NoError(t, err)

	c := &Container{
		Details: container.InspectResponse{
			ContainerJSONBase: &container.ContainerJSONBase{},
			Config: &container.Config{
				ExposedPorts: nat.PortSet{
					port: struct{}{},
				},
			},
			NetworkSettings: &container.NetworkSettings{
				Networks: map[string]*network.EndpointSettings{
					"bridge": {IPAddress: "172.17.0.2"},
				},
			},
		},
	}

	output, err := c.RenderPorts(tr)
	assert.NoError(t, err)
	assert.Contains(t, output, "443/tcp")
	assert.Contains(t, output, "-")
}

func TestRenderPortsNothingToDisplay(t *testing.T) {
	tr := testTranslationSet(t)
	c := &Container{
		Details: container.InspectResponse{
			ContainerJSONBase: &container.ContainerJSONBase{},
			Config:            &container.Config{},
			NetworkSettings:   &container.NetworkSettings{},
		},
	}

	output, err := c.RenderPorts(tr)
	assert.NoError(t, err)
	assert.Equal(t, tr.NothingToDisplay, output)
}
