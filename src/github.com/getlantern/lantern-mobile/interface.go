package client

import (
	"github.com/getlantern/balancer"
	"github.com/getlantern/lantern-mobile/interceptor"
)

var (
	i *interceptor.Interceptor
)

// Getfiretweetversion returns the current build version string
func GetFireTweetVersion() string {
	return defaultClient.getFireTweetVersion()
}

// GoCallback is the supertype of callbacks passed to Go
type GoCallback interface {
	AfterConfigure()
	AfterStart()
	WritePacket([]byte)
}

type SocketProvider interface {
	Protect(fileDescriptor int) error
}

// RunClientProxy creates a new client at the given address.
func RunClientProxy(listenAddr, appName string, protector SocketProvider, ready GoCallback) error {
	go func() {
		defaultClient = newClient(listenAddr, appName, protector)
		defaultClient.serveHTTP()
		ready.AfterStart()
	}()
	return nil
}

func IsMasqueradeCheck(ip string) bool {
	return defaultClient.IsMasqueradeCheck(ip)
}

func getBalancer() *balancer.Balancer {
	return defaultClient.Client.GetBalancer()
}

func Configure(protector SocketProvider, httpAddr string,
	socksAddr string, ready GoCallback) error {
	go func() {
		balancer.Protector = protector
		i = interceptor.New(protector, false,
			httpAddr,
			socksAddr,
			//getBalancer(),
			ready.WritePacket, IsMasqueradeCheck)
		ready.AfterConfigure()
	}()
	return nil
}

// StopClientProxy stops the proxy.
func StopClientProxy() error {
	defaultClient.stop()
	return nil
}
