package collect

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/karantan/gofp"
)

// Comparable interface for CPU_Usage
type Comparable interface {
	EqualTo(j CPU_Usage) bool
}

// EqualTo compares CPU_Usage by `.PID`
func (i CPU_Usage) EqualTo(j CPU_Usage) bool {
	return i.PID == j.PID
}

// CPU_Usage holds information needed to calculate CPU credits usage
type CPU_Usage struct {
	PCPU     float64 // cpu utilization of the process in %
	Duration float64 // cumulative CPU time in seconds
	Program  string  // name of the (parent) program
	PID      int     // pid of the process
	User     string  // effective user name
}

// define psStats function for DI
type psStats func([]string) []CPU_Usage

// ProgramCPU add/updates cpu usage for all `program` processes (one or many).
// Return 2 slices of CPU_Usage. First slice contain processes that are still running
// and the 2nd slice contains processes that have stopped. These can be used for
// further processing (e.g. send it to a DB, or just print it to stdout).
//
// Algorithm:
// 1. gather process stats for the program name `program`
// 2. loop through the given `oldUsage` CPU_Usage slice and check if every element in the
//    `oldUsage` exists in the `newUsage` (that we got from step 1)
// 2.1 Those that are not found in the `newUsage` add to the `stopped` CPU_Usage slice
// 3. loop through the `newUsage` and add/update CPU_Usage elements in the `running`
// 	  CPU_Usage slice.
func ProgramCPU(programs []string, oldUsage []CPU_Usage, fn psStats) (running []CPU_Usage, stopped []CPU_Usage) {
	newUsage := fn(programs)

	// step 2 (filter out processes that have stopped)
	notNewMember := func(i CPU_Usage) bool {
		return !Member(i, newUsage)
	}
	stopped = gofp.Filter(notNewMember, oldUsage)

	// step 3 (add and update running processes)
	isOldMember := func(i CPU_Usage) bool {
		return Member(i, oldUsage)
	}
	existing := gofp.Filter(isOldMember, newUsage)
	running = append(running, existing...)

	notOldMember := func(i CPU_Usage) bool {
		return !Member(i, oldUsage)
	}
	new := gofp.Filter(notOldMember, newUsage)
	running = append(running, new...)

	return
}

// GetProgramStats is a wrapper for calling `parseStatPS(statFromPS(<program>))`
func GetProgramStats(programs []string) (usages []CPU_Usage) {
	return parseStatPS(statFromPS(programs))
}

func parseStatPS(psOut string) (usages []CPU_Usage) {
	psOut = strings.TrimSpace(psOut)
	lines := strings.Split(psOut, "\n")

	for _, line := range lines {
		infoArr := strings.Fields(line)
		usages = append(usages, CPU_Usage{
			PCPU:     parseFloat(infoArr[0]),
			Duration: parseCPUTime(infoArr[1]) + 0.5, // because we are checking every second
			Program:  infoArr[4],
			PID:      parseInt(infoArr[2]),
			User:     infoArr[3],
		})
	}
	return
}

// statFromPS returns a programs (average) CPU usage since it started running. If the
// program has forks it will also return them (e.g. php-fpm child workers).
// Example output:
// `
//  0.1 00:00:05 3477510 root myprogram
//  0.4 00:00:01 3518860 worker1 myprogram
// 16.1 00:00:03 3519918 worker2 myprogram
//  8.0 00:00:12 3520027 worker3 myprogram
// `
// The first number is %CPU usage, 2nd is the cumulative CPU time, process pid and
// effective user name.
// For more details see https://man7.org/linux/man-pages/man1/ps.1.html
func statFromPS(programs []string) string {
	f := func(p string) string {
		return fmt.Sprintf("$(pgrep %s)", p)
	}
	pids := strings.Join(gofp.ForEach(f, programs), " ")

	psCommand := fmt.Sprintf("ps -o pcpu=,time=,pid=,user:32=,comm= %s", pids)
	cmd := exec.Command("bash", "-c", psCommand)
	log.Debugf("Running: `%s`", cmd.String())
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()

	if err != nil {
		log.Error(out)
		log.Error(err)
	}

	return string(out)
}

// CPUs returns the number of vCPUs the system has
func CPUs() int {
	// `nproc` is part of coreutils. See https://man7.org/linux/man-pages/man1/nproc.1.html
	cmd := exec.Command("nproc")
	log.Debugf("Running: `%s`", cmd.String())
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()

	if err != nil {
		log.Error(out)
		log.Error(err)
	}
	cpus, err := strconv.Atoi(string(out))
	if err != nil {
		log.Error(err)
	}

	return cpus
}

//
// Util functions
//

func parseInt(val string) int {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Error(err)
	}
	return intVal
}

func parseFloat(val string) float64 {
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Error(err)
	}
	return floatVal
}

// parseCPUTime parses cumulative CPU time, "[DD-]hh:mm:ss" format to int (seconds only).
func parseCPUTime(val string) (duration float64) {
	if strings.Contains(val, "-") {
		splitDays := strings.Split(val, "-")
		days := splitDays[0]

		if days != "" {
			secInDay := 24 * 60 * 60
			duration += float64(secInDay * parseInt(days)) // add days
		}
		val = splitDays[1]
	}

	splitTime := strings.Split(val, ":")
	hours := parseFloat(splitTime[0])
	duration += hours * 60 * 60 // add hours
	min := parseFloat(splitTime[1])
	duration += min * 60                 // add minutes
	duration += parseFloat(splitTime[2]) // add seconds

	return
}

// Member checks if an `element` exists in the given `slice`. Returns true otherwise false.
// See https://package.elm-lang.org/packages/elm/core/latest/List#member
func Member(element CPU_Usage, slice []CPU_Usage) bool {
	for _, v := range slice {
		if element.PID == v.PID {
			return true
		}
	}
	return false
}
