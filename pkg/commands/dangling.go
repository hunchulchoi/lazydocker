package commands

// IsDangling reports whether the image matches `docker images -f dangling=true`.
// Dangling images are untagged and not used by any container.
func (i *Image) IsDangling() bool {
	if i.Image.Containers != 0 {
		return false
	}

	return len(i.Image.RepoTags) == 0
}

// IsDangling reports whether the volume matches `docker volume ls -f dangling=true`.
func (v *Volume) IsDangling() bool {
	if v.Volume == nil || v.Volume.UsageData == nil {
		return false
	}

	return v.Volume.UsageData.RefCount == 0
}

// IsDangling reports whether the network matches `docker network ls -f dangling=true`.
func (n *Network) IsDangling() bool {
	return len(n.Network.Containers) == 0
}
