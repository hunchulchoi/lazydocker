package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageIsDangling(t *testing.T) {
	assert.True(t, (&Image{Dangling: true}).IsDangling())
	assert.False(t, (&Image{Dangling: false}).IsDangling())
}

func TestVolumeIsDangling(t *testing.T) {
	assert.True(t, (&Volume{Dangling: true}).IsDangling())
	assert.False(t, (&Volume{Dangling: false}).IsDangling())
}

func TestNetworkIsDangling(t *testing.T) {
	assert.True(t, (&Network{Dangling: true}).IsDangling())
	assert.False(t, (&Network{Dangling: false}).IsDangling())
}

func TestMapContains(t *testing.T) {
	set := map[string]struct{}{"a": {}}
	assert.True(t, mapContains(set, "a"))
	assert.False(t, mapContains(set, "b"))
	assert.False(t, mapContains(nil, "a"))
}
