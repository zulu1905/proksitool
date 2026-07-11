//go:build !linux

package dnscrypt

func (proxy *Proxy) addSystemDListeners() error {
	return nil
}
