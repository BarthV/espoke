// Copyright © 2018 Barthelemy Vessemont
// GNU General Public License version 3

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var consulTarget string
var consulPeriod string
var probePeriod string
var cleanMetricsPeriod string
var metricsPort int
var loglevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "espoke",
	Short: "espoke is a whitebox probing tool for Elasticsearch clusters",
	Long: `espoke is a whitebox probing tool for Elasticsearch clusters.
It completes the following actions :
 * discover every ES clusters registered in Consul
 * run an empty search query against every discovered indexes, data servers & clusters
 * expose latency metrics with tags for clusters and nodes
 * expose avaibility metrics with tags for clusters and nodes`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.espoke.yaml)")
	rootCmd.PersistentFlags().StringVarP(&consulTarget, "consulApi", "a", "127.0.0.1:8500", "consul target api host:port")
	rootCmd.PersistentFlags().StringVar(&consulPeriod, "consulPeriod", "120s", "nodes discovery update interval")
	rootCmd.PersistentFlags().StringVar(&cleanMetricsPeriod, "cleaningPeriod", "600s", "prometheus metrics cleaning interval (for vanished nodes)")
	rootCmd.PersistentFlags().StringVar(&probePeriod, "probePeriod", "30s", "elasticsearch nodes probing interval")
	rootCmd.PersistentFlags().IntVarP(&metricsPort, "metricsPort", "p", 2112, "port where prometheus will expose metrics to")
	rootCmd.PersistentFlags().StringVarP(&loglevel, "loglevel", "l", "info", "log level")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".espoke" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".espoke")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Init logger
	log.SetOutput(os.Stdout)
	lvl, err := log.ParseLevel(loglevel)
	if err != nil {
		log.Warning("Log level not recognized, fallback to default level (INFO)")
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)
	log.Info("Logger initialized")
}
