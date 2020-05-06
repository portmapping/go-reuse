package reuse

import (
	"net"
)

var addrMapping = map[string]func(network, address string) (net.Addr, error){
	"ip":         resolveIPAddr,
	"ip4":        resolveIPAddr,
	"ip6":        resolveIPAddr,
	"tcp":        resolveTCPAddr,
	"tcp4":       resolveTCPAddr,
	"tcp6":       resolveTCPAddr,
	"udp":        resolveUDPAddr,
	"udp4":       resolveUDPAddr,
	"udp6":       resolveUDPAddr,
	"unix":       resolveUnixAddr,
	"unixgram":   resolveUnixAddr,
	"unixpacket": resolveUnixAddr,
}

func ResolveAddr(network, address string) (net.Addr, error) {
	if v, b := addrMapping[network]; b {
		return v(network, address)
	}
	return nil, net.UnknownNetworkError(network)

}

func resolveIPAddr(network, address string) (net.Addr, error) {
	return net.ResolveIPAddr(network, address)
}

func resolveTCPAddr(network, address string) (net.Addr, error) {
	return net.ResolveTCPAddr(network, address)
}

func resolveUDPAddr(network, address string) (net.Addr, error) {
	return net.ResolveUDPAddr(network, address)
}

func resolveUnixAddr(network, address string) (net.Addr, error) {
	return net.ResolveUnixAddr(network, address)
}
