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
	"time"
)

var (
	log      = logger.New("main")
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
	initInflux(os.Getenv("INFLUXDB_PSUSAGE_HOST"))
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
	log.Infof("Running psusage version %s", version.Version)
	var program string
	flag.StringVar(&program, "program", "", "Name of the program you want to monitor CPU usage over time. Example: mysqld, haproxy, php-fpm, etc")
	flag.Parse()

	log.Infof("Monitoring CPU usage for %s", program)

	c := time.Tick(1 * time.Second)
	running := []collect.CPU_Usage{}
	stopped := []collect.CPU_Usage{}
	for ; true; <-c {
		running, stopped = collect.ProgramCPU(program, running, collect.GetProgramStats)
		if len(stopped) > 0 {
			for _, p := range stopped {
				influxdb.AddPoint(influx, p, hostname)
				log.Infof("%s (%s:%d) used %f%% CPU over %d seconds.", p.Program, p.User, p.PID, p.PCPU, p.Duration)
			}
			stopped = []collect.CPU_Usage{}
		}
	}
}
