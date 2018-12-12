// Copyright © 2018 Barthelemy Vessemont
// GNU General Public License version 3

package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
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

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr == nil {
		var p fastjson.Parser
		v, jsonErr := p.Parse(string(body))
		if jsonErr == nil {
			shardsSuccessfulGauge.WithLabelValues(node.cluster, node.name).Set(v.GetFloat64("_shards", "successful"))
			docsHitGauge.WithLabelValues(node.cluster, node.name).Set(v.GetFloat64("hits", "total"))
		}
	}

	nodeAvailabilityGauge.WithLabelValues(node.cluster, node.name).Set(1)
	nodeSearchLatencySummary.WithLabelValues(node.cluster, node.name).Observe(durationNanosec)

	return nil
}
