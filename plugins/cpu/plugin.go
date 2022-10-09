package cpu

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/spf13/cobra"
)

type CPU struct {
	cpuUsage prometheus.Gauge
}

func (c *CPU) Init(cmd *cobra.Command) error {
	c.cpuUsage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_all",
		Help: "The cpu usage percentage, all core(s)",
	})
	return nil
}

func (c *CPU) Update() error {
	value, err := cpu.Percent(time.Second*1, false)
	if err != nil {
		return err
	}
	c.cpuUsage.Set(value[0])
	return nil
}

func (c *CPU) GetName() string {
	return "cpu"
}
