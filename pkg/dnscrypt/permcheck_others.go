//go:build !unix

package dnscrypt

func WarnIfMaybeWritableByOtherUsers(p string) {
	// No-op
}
