//go:build !freebsd && !openbsd && !windows && !darwin && !linux

package dnscrypt

import (
	"net"
)

func (proxy *Proxy) udpListenerConfig() (*net.ListenConfig, error) {
	return &net.ListenConfig{}, nil
}

func (proxy *Proxy) tcpListenerConfig() (*net.ListenConfig, error) {
	return &net.ListenConfig{}, nil
}
