package hddtemp

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/cobra"
)

type HDDTemp struct {
	temp *prometheus.GaugeVec
	addr string
}

func (h *HDDTemp) Init(cmd *cobra.Command) error {
	h.temp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hdd_temp",
			Help: "HDD Temperature",
		},
		[]string{"name"},
	)
	cmd.Flags().StringVar(&h.addr, "hddtemp_addr", "localhost:7634", "hddtemp server address")
	return nil
}

func (h *HDDTemp) GetName() string {
	return "hdd_temperature_sensors"
}

func (h *HDDTemp) Update() error {
	var (
		err    error
		conn   net.Conn
		buffer bytes.Buffer
	)

	if conn, err = net.Dial("tcp", h.addr); err != nil {
		log.Println(err)
		return err
	}

	if _, err = io.Copy(&buffer, conn); err != nil {
		log.Println(err)
		return err
	}

	fields := strings.Split(buffer.String(), "|")

	for index := 0; index < len(fields)/5; index++ {
		offset := index * 5
		device := fields[offset+1]
		device = device[strings.LastIndex(device, "/")+1:]

		temperatureField := fields[offset+3]
		temperature, err := strconv.ParseInt(temperatureField, 10, 32)

		if err != nil {
			log.Println("error while parse temperature, set to 0")
			temperature = 0
		}
		h.temp.With(prometheus.Labels{"name": device}).Set(float64(temperature))
	}

	return nil
}
