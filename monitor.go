package main

import (
	"fmt"
	"strings"
	"time"

	"net/http"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/ssimunic/gosensors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_all",
		Help: "The cpu usage percentage, all core(s)",
	})
	memTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mem_total",
		Help: "Total memory",
	})
	memFree = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mem_free",
		Help: "Free memory",
	})
	memUsed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mem_used",
		Help: "Total memory",
	})
)

func updateMetrics() {
	go func() {
		for {
			v, _ := mem.VirtualMemory()
			c, _ := cpu.Percent(time.Second*1, false)

			cpuUsage.Set(c[0])
			memTotal.Set(float64(v.Total))
			memFree.Set(float64(v.Free))
			memUsed.Set(v.UsedPercent)

			sensors, _ := gosensors.NewFromSystem()
			for chip := range sensors.Chips {
				// Iterate over entries
				for key, value := range sensors.Chips[chip] {
					// If CPU or GPU, print out
					if strings.HasPrefix(key, "Core") {
						fmt.Println(key, value)
					}
				}
			}
			fmt.Println(sensors)
			time.Sleep(time.Second * 10)
		}
	}()
}

func main() {
	updateMetrics()
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
