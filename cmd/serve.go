// Copyright © 2018 Barthelemy Vessemont
// GNU General Public License version 3

package cmd

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP endpoint, ES discovery and probing",
	Long: `Start ES discovering mecanism & probe periodically requests against all ES nodes.
Expose all measures using a prometheus compliant HTTP endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Entering serve main loop")
		startMetricsEndpoint()

		log.Info("Discovering ES nodes for the first time")
		allEverKnownDatanodes := []string{}
		datanodesList, err := discoverEsNodes()
		if err != nil {
			errorsCount.Inc()
			log.Fatal("Impossible to discover ES datanodes during bootstrap, exiting")
		}
		allEverKnownDatanodes = updateEverKnownDatanodes(allEverKnownDatanodes, datanodesList)

		log.Info("Initializing tickers")
		updateDiscoveryPeriod, err := time.ParseDuration(consulPeriod)
		if err != nil {
			log.Warning("Impossible to parse consulPeriod value, fallback to 120s")
			updateDiscoveryPeriod = 120 * time.Second
		}
		if updateDiscoveryPeriod < 60*time.Second {
			log.Warning("Refreshing discovery more than once a minute is not allowed, fallback to 60s")
			updateDiscoveryPeriod = 60 * time.Second
		}
		log.Info("Discovery update interval: ", updateDiscoveryPeriod.String())

		updateProbingPeriod, err := time.ParseDuration(probePeriod)
		if err != nil {
			log.Warning("Impossible to parse probePeriod value, fallback to 30s")
			updateProbingPeriod = 30 * time.Second
		}
		if updateProbingPeriod < 20*time.Second {
			log.Warning("Probing elasticsearch datanodes more than 3 times a minute is not allowed, fallback to 20s")
			updateProbingPeriod = 20 * time.Second
		}
		log.Info("Probing interval: ", updateProbingPeriod.String())

		pruneMetricsPeriod, err := time.ParseDuration(cleanMetricsPeriod)
		if err != nil {
			log.Warning("Impossible to parse cleaningPeriod value, fallback to 600s")
			pruneMetricsPeriod = 600 * time.Second
		}
		if pruneMetricsPeriod < 240*time.Second {
			log.Warning("Cleaning Metrics faster than every 4 minutes is not allowed, fallback to 240s")
			pruneMetricsPeriod = 240 * time.Second
		}
		log.Info("Metrics pruning interval: ", pruneMetricsPeriod.String())

		updateDiscoveryTicker := time.NewTicker(updateDiscoveryPeriod)
		cleanMetricsTicker := time.NewTicker(pruneMetricsPeriod)
		executeProbingTicker := time.NewTicker(updateProbingPeriod)

		for {
			select {
			case <-cleanMetricsTicker.C:
				log.Info("Cleaning Prometheus metrics for unreferenced nodes")
				cleanMetrics(datanodesList, allEverKnownDatanodes)

			case <-updateDiscoveryTicker.C:
				log.Debug("Starting updating ES data nodes list")

				updatedList, err := discoverEsNodes()
				if err != nil {
					log.Error("Unable to update ES datanodes, using last known state")
					errorsCount.Inc()
					continue
				}

				log.Info("Updating data nodes list")
				allEverKnownDatanodes = updateEverKnownDatanodes(allEverKnownDatanodes, updatedList)
				datanodesList = updatedList

			case <-executeProbingTicker.C:
				log.Debug("Starting probing ES nodes")

				sem := new(sync.WaitGroup)
				for _, node := range datanodesList {
					sem.Add(1)
					go func(loopnode datanode) {
						defer sem.Done()
						probeDatanode(&loopnode, updateProbingPeriod)
					}(node)

				}
				sem.Wait()

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
