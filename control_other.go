// +build !windows,!linux,!darwin,!dragonfly,!freebsd,!netbsd,!openbsd

package reuse

import (
	"syscall"
)

func Control(network, address string, c syscall.RawConn) error {
	return nil
}
