package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	avgProc = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "avgproc",
			Help: "AVGPROC item from vmcp indicate",
		},
		func() float64 {
			return getAvgProc()
		},
	)
)

func runVMCP() string {
	cmd := exec.Command("/home/chale/zvm_exporter/test_data.sh")

	var out bytes.Buffer

	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return out.String()

}

func getAvgProc() float64 {
	vmcpInput := runVMCP()

	avgproc := strings.TrimLeft(vmcpInput, "AVGPROC-")

	retVal, err := strconv.ParseFloat(avgproc[0:3], 32)
	if err != nil {
		os.Exit(1)
	}

	return retVal
}

func main() {
	fmt.Println("Starting z/VM Exporter!")

	path := "/metrics"
	addr := ":9800"

	metricsPath := &path

	prometheus.MustRegister(avgProc)

	//level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	http.Handle(*metricsPath, promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Haproxy Exporter</title></head>
             <body>
             <h1>Haproxy Exporter</h1>
             <p><a href='` + path + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}
}