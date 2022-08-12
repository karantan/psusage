package main

import (
	"flag"
	"fmt"
	"os"
	"psusage/collect"
	"psusage/influxdb"
	"psusage/logger"
	"psusage/version"
	"runtime"
	"strings"
	"time"
)

var (
	log      = logger.New("main", false)
	hostname string
	influx   influxdb.InfluxClient
)

func init() {
	platform := runtime.GOOS
	if platform != "linux" {
		fmt.Println("Only supported on Linux.")
		os.Exit(1)
	}
	hostname, _ = os.Hostname()

	influxDSN := os.Getenv("INFLUXDB_PSUSAGE_DSN")
	if influxDSN == "" {
		fmt.Println("INFLUXDB_PSUSAGE_DSN env var is missing.")
		os.Exit(1)
	}
	initInflux(influxDSN)
}

func initInflux(influxDSN string) influxdb.InfluxClient {
	influx = influxdb.NewInfluxDSN(influxDSN)
	_, _, err := influx.Ping(5)
	if err != nil {
		log.Fatal("No response from the InfluxDB: ", err)
		os.Exit(1)
	}
	return influx
}

func main() {
	var p string
	flag.StringVar(&p, "programs", "", `Name of the program(s) you want to monitor CPU usage over time. Example:

	psusage --programs "mysqld haproxy php-fpm"

Or just:

	psusage --programs mysqld`,
	)
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Debug mode. I'll write the data to the InfluxDB and also to the stdout.")
	flag.Parse()
	programs := strings.Fields(p)

	log.Infof("Running psusage version %s", version.Version)
	log.Infof("Monitoring CPU usage for %s", programs)

	c := time.Tick(1 * time.Second)
	running := []collect.CPU_Usage{}
	stopped := []collect.CPU_Usage{}
	for ; true; <-c {
		running, stopped = collect.ProgramCPU(programs, running, collect.GetProgramStats)
		if len(stopped) > 0 {
			for _, p := range stopped {
				influxdb.AddPoint(influx, p, hostname)

				if debug {
					log.Infof("%s (%s:%d) used %f%% CPU over %d seconds.", p.Program, p.User, p.PID, p.PCPU, p.Duration)
				}
			}
			stopped = []collect.CPU_Usage{}
		}
	}
}
