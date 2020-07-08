package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type zVMStats struct {
	avgProc float64
}

var (
	pollInterval int
	vmcpPath     *string
	zvmStats     zVMStats
	avgProc      = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "avgproc",
			Help: "AVGPROC item from vmcp indicate",
		},
		func() float64 {
			return zvmStats.avgProc
		},
	)
)

func runVMCP() string {
	cmd := exec.Command(*vmcpPath + " " + "indicate")

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

func parseVMCP(stats *zVMStats) {
	stats.avgProc = getAvgProc()
}

func updateLoop() {

	for {
		parseVMCP(&zvmStats)
		time.Sleep(time.Duration(pollInterval))
	}
}

func main() {

	metricsPath := flag.String("path", "/metrics", "URL Path for metrics endpoint")
	addr := flag.String("address", ":9100", "Bind address")
	pollIntervalPtr := flag.Int("pollInterval", 15, "Interval to poll z/VM")
	vmcpPath = flag.String("vmcpPath", "/home/chale/zvm_exporter/test_data.sh", "Path to VMCP Utility")

	flag.Parse()

	pollInterval = *pollIntervalPtr

	fmt.Println("Starting z/VM Exporter!")
	fmt.Println("Listening on " + *addr + *metricsPath)

	go updateLoop()

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
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	if err := http.ListenAndServe(*addr, nil); err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}
}
