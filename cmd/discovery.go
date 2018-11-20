// Copyright © 2018 Barthelemy Vessemont
// GNU General Public License version 3

package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

type datanode struct {
	name    string
	ip      string
	port    int
	cluster string
	version string
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func updateEverKnownDatanodes(allEverKnownDatanodes []string, datanodes []datanode) []string {
	for _, node := range datanodes {
		serializedNode := fmt.Sprintf("%v|%v", node.name, node.cluster)
		if contains(allEverKnownDatanodes, serializedNode) == false {
			allEverKnownDatanodes = append(allEverKnownDatanodes, serializedNode)
		}
	}
	sort.Strings(allEverKnownDatanodes)
	return allEverKnownDatanodes
}

func clusterNameFromTags(serviceTags []string) string {
	for _, tag := range serviceTags {
		splitted := strings.SplitN(tag, "-", 2)
		if splitted[0] == "cluster_name" {
			return splitted[1]
		}
	}
	return ""
}

func versionFromTags(serviceTags []string) string {
	for _, tag := range serviceTags {
		splitted := strings.SplitN(tag, "-", 2)
		if splitted[0] == "version" {
			return splitted[1]
		}
	}
	return ""
}

func discoverEsNodes() ([]datanode, error) {
	start := time.Now()

	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulTarget
	consul, err := api.NewClient(consulConfig)
	if err != nil {
		log.Debug("Consul Connection failed: ", err.Error())
		errorsCount.Inc()
		return nil, err
	}

	catalogServices, _, err := consul.Catalog().ServiceMultipleTags(
		"elasticsearch-all", []string{"data"},
		&api.QueryOptions{AllowStale: true, RequireConsistent: false, UseCache: true},
	)
	if err != nil {
		log.Error("Consul Discovery failed: ", err.Error())
		errorsCount.Inc()
		return nil, err
	}

	var datanodeList []datanode
	for _, svc := range catalogServices {
		log.Debug("Service discovered: ", svc.Node, " (", svc.Address, ":", svc.ServicePort, ")")

		datanodeList = append(datanodeList, datanode{
			name:    svc.Node,
			ip:      svc.Address,
			port:    svc.ServicePort,
			cluster: clusterNameFromTags(svc.ServiceTags),
			version: versionFromTags(svc.ServiceTags),
		})
	}

	nodesCount := len(datanodeList)
	datanodeCount.Set(float64(nodesCount))
	log.Debug(nodesCount, " datanodes found")

	consulDiscoveryDurationSummary.Observe(float64(time.Since(start).Nanoseconds()))
	return datanodeList, nil
}
