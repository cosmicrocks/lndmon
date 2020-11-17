package collectors

import (
	"context"
	"encoding/hex"

	"github.com/lightninglabs/lndclient"
	"github.com/prometheus/client_golang/prometheus"
)

// InfoCollector is a collector that keeps track of node information.
type InfoCollector struct {
	info *prometheus.Desc

	lnd lndclient.LightningClient
}

// NewInfoCollector returns a new instance of the InfoCollector for the target
// lnd client.
func NewInfoCollector(lnd lndclient.LightningClient) *InfoCollector {
	labels := []string{"version", "alias", "pubkey"}
	return &InfoCollector{
		info: prometheus.NewDesc(
			"lnd_info", "lnd node info", labels, nil,
		),
		lnd: lnd,
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once the
// last descriptor has been sent.
//
// NOTE: Part of the prometheus.Collector interface.
func (c *InfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.info
}

// Collect is called by the Prometheus registry when collecting metrics.
//
// NOTE: Part of the prometheus.Collector interface.
func (c *InfoCollector) Collect(ch chan<- prometheus.Metric) {
	resp, err := c.lnd.GetInfo(context.Background())
	if err != nil {
		chainLogger.Error(err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.info, prometheus.GaugeValue, 0, resp.Version,
		resp.Alias, hex.EncodeToString(resp.IdentityPubkey[:]),
	)
}
