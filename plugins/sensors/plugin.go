package sensors

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/cobra"
	"github.com/ssimunic/gosensors"
)

type Sensors struct {
	temp   *prometheus.GaugeVec
	regexp *regexp.Regexp
}

func (s *Sensors) Init(cmd *cobra.Command) error {
	s.temp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temp",
			Help: "Temperature",
		},
		[]string{"name"},
	)
	s.regexp = regexp.MustCompile(`^\+([0-9.]+)°C.*$`)

	return nil
}

func (s *Sensors) Update() error {
	sensors, _ := gosensors.NewFromSystem()
	for chip := range sensors.Chips {
		for key, value := range sensors.Chips[chip] {
			name := fmt.Sprintf("%s-%s", chip, strings.ReplaceAll(key, " ", "-"))
			if s.regexp.MatchString(value) {
				match := s.regexp.FindStringSubmatch(value)
				temperature, err := strconv.ParseFloat(match[1], 64)
				if err == nil {
					s.temp.With(prometheus.Labels{"name": name}).Set(temperature)
				}
			}
		}
	}
	return nil
}

func (s *Sensors) GetName() string {
	return "temperature_sensors"
}
