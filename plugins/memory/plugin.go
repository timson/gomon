package memory

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

type Memory struct {
	memTotal prometheus.Gauge
	memFree  prometheus.Gauge
	memUsed  prometheus.Gauge
}

func (m *Memory) Init(cmd *cobra.Command) error {
	m.memTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mem_total",
		Help: "Total memory",
	})
	m.memFree = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mem_free",
		Help: "Free memory",
	})
	m.memUsed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mem_used",
		Help: "Total memory",
	})

	return nil
}

func (m *Memory) Update() error {
	value, err := mem.VirtualMemory()
	if err != nil {
		return nil
	}
	m.memTotal.Set(float64(value.Total))
	m.memFree.Set(float64(value.Free))
	m.memUsed.Set(value.UsedPercent)
	return nil
}

func (m *Memory) GetName() string {
	return "memory"
}
