// Copyright © 2018 Barthelemy Vessemont
// GNU General Public License version 3

package cmd

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func probeNode(node *esnode, updateProbingPeriod time.Duration) error {
	client := &http.Client{
		Timeout: updateProbingPeriod - 2*time.Second,
	}

	probingURL := fmt.Sprintf("http://%v:%v/_all/_search", node.ip, node.port)
	log.Debug("Start probing ", node.name)

	start := time.Now()
	resp, err := client.Get(probingURL)
	if err != nil {
		log.Debug("Probing failed for ", node.name, ": ", probingURL, " ", err.Error())
		log.Error(err)
		nodeAvailabilityGauge.WithLabelValues(node.cluster, node.name).Set(0)
		errorsCount.Inc()
		return err
	}
	durationNanosec := float64(time.Since(start).Nanoseconds())

	log.Debug("Probe result for ", node.name, ": ", resp.Status)
	if resp.StatusCode != 200 {
		log.Error("Probing failed for ", node.name, ": ", probingURL, " ", resp.Status)
		nodeAvailabilityGauge.WithLabelValues(node.cluster, node.name).Set(0)
		errorsCount.Inc()
		return fmt.Errorf("ES Probing failed")
	}

	nodeAvailabilityGauge.WithLabelValues(node.cluster, node.name).Set(1)
	nodeSearchLatencySummary.WithLabelValues(node.cluster, node.name).Observe(durationNanosec)

	return nil
}
