package main

import (
	"flag"
	"psusage/collect"
	"psusage/logger"
	"psusage/version"
	"time"
)

// func init() {
// 	platform := runtime.GOOS
// 	if platform != "linux" {
// 		fmt.Println("Only supported on Linux.")
// 		os.Exit(1)
// 	}
// }

var log = logger.New("main")

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
				log.Infof("%s (%s:%d) used %f%% CPU over %d seconds.", p.Program, p.User, p.PID, p.PCPU, p.Duration)
			}
			stopped = []collect.CPU_Usage{}
		}
	}
}
