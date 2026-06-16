package commands

// IsDangling reports whether the image matches `docker images -f dangling=true`.
func (i *Image) IsDangling() bool {
	return i.Dangling
}

// IsDangling reports whether the volume matches `docker volume ls -f dangling=true`.
func (v *Volume) IsDangling() bool {
	return v.Dangling
}

// IsDangling reports whether the network matches `docker network ls -f dangling=true`.
func (n *Network) IsDangling() bool {
	return n.Dangling
}
