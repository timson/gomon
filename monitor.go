package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/timson/gomon/plugins"
	"github.com/timson/gomon/plugins/cpu"
	"github.com/timson/gomon/plugins/hddtemp"
	"github.com/timson/gomon/plugins/memory"
	"github.com/timson/gomon/plugins/sensors"
)

var version = "0.0.1"
var rootCmd = &cobra.Command{
	Use:     "gomon",
	Version: version,
	Short:   "gomon - a simple monitoring tool with export to Prometheus",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var port int
var interval int
var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve prometheus metrics",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		serve(port, interval)
	},
}
var activePlugins = []plugins.Plugin{
	&cpu.CPU{},
	&memory.Memory{},
	&sensors.Sensors{},
	&hddtemp.HDDTemp{},
}

func serve(port int, interval int) {
	go func() {
		for {
			for _, plugin := range activePlugins {
				err := plugin.Update()
				if err != nil {
					log.Println(err)
				}
			}
			time.Sleep(time.Second * time.Duration(interval))
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("start serve /metrics endpoint on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	for _, plugin := range activePlugins {
		err := plugin.Init(serverCmd)
		if err == nil {
			log.Printf("plugin %s loaded", plugin.GetName())
		} else {
			log.Println(err)
		}
	}
	serverCmd.Flags().IntVar(&port, "port", 8080, "port to serve")
	serverCmd.Flags().IntVar(&interval, "interval", 10, "metrics update interval in seconds")
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
