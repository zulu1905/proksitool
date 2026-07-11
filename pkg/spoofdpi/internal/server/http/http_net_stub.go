//go:build !darwin

package http

import (
	"github.com/rs/zerolog"
	"myvpn/pkg/spoofdpi/internal/netutil"
	"myvpn/pkg/spoofdpi/internal/packet"
)

const StateFile = ""

type httpSystemNetworkStub struct{}

func NewHTTPSystemNetwork(
	logger zerolog.Logger,
	defaultRoute *packet.Route,
) HTTPSystemNetwork {
	return &httpSystemNetworkStub{}
}

func (n *httpSystemNetworkStub) DefaultRoute() *packet.Route {
	return nil
}

func (n *httpSystemNetworkStub) BuildJobs(
	_ uint16,
	_ string,
) ([]netutil.NetworkJob, error) {
	return nil, nil
}
