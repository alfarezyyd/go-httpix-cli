package outbound

import (
	"sync"
)

type Outbound struct {
	HttpClient *HttpClient
}

var (
	outboundInstance *Outbound
	syncOnce         sync.Once
)

func Setup() *Outbound {
	httpClient := httpClientSetup()

	return &Outbound{HttpClient: httpClient}
}

func GetHttpClient() *HttpClient {
	syncOnce.Do(func() {
		outboundInstance = Setup()
	})
	return outboundInstance.HttpClient
}
