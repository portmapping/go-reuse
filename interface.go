// Package reuse provides Listen and Dial functions that set socket
// options in order to be able to reuse ports. You should only use this
// package if you know what SO_REUSEADDR and SO_REUSEPORT are.
//
// For example:
//
//  // listen on the same port.
//  l1, _ := greuse.Listen("tcp", "127.0.0.1:1234")
//  l2, _ := greuse.Listen("tcp", "127.0.0.1:1234")
//
//  // dial from the same port.
//  l1, _ := greuse.Listen("tcp", "127.0.0.1:1234")
//  l2, _ := greuse.Listen("tcp", "127.0.0.1:1235")
//  c, _ := greuse.Dial("tcp", "127.0.0.1:1234", "127.0.0.1:1235")
//
// Note: can't dial self because tcp/ip stacks use 4-tuples to identify connections,
// and doing so would clash.
package reuse

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

var (
	// Enabled returns whether or not SO_REUSEPORT or equivalent behaviour is
	// enabled in the OS.
	Enabled      = false
	listenConfig = net.ListenConfig{
		Control: Control,
	}
)

// Listen listens at the given network and address. see net.Listen
// Returns a net.Listener created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func Listen(network, address string) (net.Listener, error) {
	return listenConfig.Listen(context.Background(), network, address)
}

// ListenTLS listens at the given network and address. see net.Listen
// Returns a net.Listener created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func ListenTLS(network, address string, config *tls.Config) (net.Listener, error) {
	listen, err := listenConfig.Listen(context.Background(), network, address)
	if err != nil {
		return nil, err
	}
	return tls.NewListener(listen, config), nil
}

// ListenTCP listens at the given network and address. see net.Listen
// Returns a net.Listener created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func ListenTCP(network string, laddr *net.TCPAddr) (*net.TCPListener, error) {
	t, err := net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}
	conn, err := t.SyscallConn()
	if err != nil {
		return nil, err
	}
	err = Control(network, "", conn)
	return t, err
}

// ListenIP listens at the given network and address. see net.Listen
// Returns a net.Listener created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func ListenIP(network string, laddr *net.IPAddr) (*net.IPConn, error) {
	i, err := net.ListenIP(network, laddr)
	if err != nil {
		return nil, err
	}
	conn, err := i.SyscallConn()
	if err != nil {
		return nil, err
	}
	err = Control(network, "", conn)
	return i, err
}

// ListenUnix listens at the given network and address. see net.Listen
// Returns a net.Listener created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func ListenUnix(network string, laddr *net.UnixAddr) (*net.UnixListener, error) {
	u, err := net.ListenUnix(network, laddr)
	if err != nil {
		return nil, err
	}
	conn, err := u.SyscallConn()
	if err != nil {
		return nil, err
	}
	err = Control(network, "", conn)
	return u, err
}

// ListenPacket listens at the given network and address. see net.ListenPacket
// Returns a net.Listener created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func ListenPacket(network, address string) (net.PacketConn, error) {
	return listenConfig.ListenPacket(context.Background(), network, address)
}

// DialTimeOut dials the given network and address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func DialTimeOut(network, laddr, raddr string, timeout time.Duration) (net.Conn, error) {
	nla, err := ResolveAddr(network, laddr)
	if err != nil {
		return nil, fmt.Errorf("resolving local addr: %w", err)
	}

	d := net.Dialer{
		Control:   Control,
		LocalAddr: nla,
		Timeout:   timeout,
	}

	return d.Dial(network, raddr)
}

// Dial dials the given network and address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func Dial(network, laddr, raddr string) (net.Conn, error) {
	nla, err := ResolveAddr(network, laddr)
	if err != nil {
		return nil, fmt.Errorf("resolving local addr: %w", err)
	}
	d := net.Dialer{
		Control:   Control,
		LocalAddr: nla,
	}
	return d.Dial(network, raddr)
}

// DialTLS dials the given network and address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func DialTLS(network, laddr, raddr string, config *tls.Config) (net.Conn, error) {
	nla, err := ResolveAddr(network, laddr)
	if err != nil {
		return nil, fmt.Errorf("resolving local addr: %w", err)
	}
	d := net.Dialer{
		Control:   Control,
		LocalAddr: nla,
	}
	return tls.DialWithDialer(&d, network, raddr, config)
}

// DialTCP dials the given network and tcp address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func DialTCP(network string, laddr *net.TCPAddr, raddr *net.TCPAddr) (net.Conn, error) {
	d := net.Dialer{
		Control:   Control,
		LocalAddr: laddr,
	}
	return d.Dial(network, raddr.String())
}

// DialAddr dials the given network and address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func DialAddr(network string, laddr net.Addr, raddr net.Addr) (net.Conn, error) {
	d := net.Dialer{
		Control:   Control,
		LocalAddr: laddr,
	}
	return d.Dial(network, raddr.String())
}

// DialIP dials the given network and ip address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func DialIP(network string, laddr *net.IPAddr, raddr *net.IPAddr) (net.Conn, error) {
	d := net.Dialer{
		Control:   Control,
		LocalAddr: laddr,
	}
	return d.Dial(network, raddr.String())
}

// DialUDP dials the given network and udp address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func DialUDP(network string, laddr *net.UDPAddr, raddr *net.UDPAddr) (net.Conn, error) {
	d := net.Dialer{
		Control:   Control,
		LocalAddr: laddr,
	}
	return d.Dial(network, raddr.String())
}

// DialTimeoutUDP dials the given network and udp address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func DialTimeoutUDP(network string, laddr *net.UDPAddr, raddr *net.UDPAddr, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control:   Control,
		LocalAddr: laddr,
		Timeout:   timeout,
	}
	return d.Dial(network, raddr.String())
}

// DialUnix dials the given network and unix address. see net.Dialer.Dial
// Returns a net.Conn created from a file discriptor for a socket
// with SO_REUSEPORT and SO_REUSEADDR option set.
func DialUnix(network string, laddr *net.UnixAddr, raddr *net.UnixAddr) (net.Conn, error) {
	d := net.Dialer{
		Control:   Control,
		LocalAddr: laddr,
	}
	return d.Dial(network, raddr.String())
}
