package influxdb

import (
	"net/url"
	"os"
	"psusage/collect"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

// InfluxClient holds InfluxDB client
type InfluxClient struct {
	client.Client
}

// NewInfluxDSN returns influxdb2 client ready for use
func NewInfluxDSN(dsn string) InfluxClient {
	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatal(err)
	}
	pass := ""
	if p, ok := u.User.Password(); ok {
		pass = p
	}
	user := u.User.Username()
	u.User = nil

	influxC, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               u.String(),
		Username:           user,
		Password:           pass,
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return InfluxClient{influxC}
}

// AddPoint adds `collect.CPU_Usage` point to the InfluxDB
func AddPoint(db InfluxClient, u collect.CPU_Usage, hostname string) {
	// Make sure the database exists (`CREATE DATABASE psusage`)
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "psusage",
		Precision: "s",
	})

	measurement := "cpu"
	tags := map[string]string{
		"server":  hostname,
		"program": u.Program,
		"user":    u.User,
	}
	fields := map[string]interface{}{
		"pcpu":     u.PCPU,
		"duration": u.Duration,
	}

	// insert cpu,server=<hostname>,program=<program>,user=<user> pcpu=<%CPU>,duration=<seconds>
	p, _ := client.NewPoint(measurement, tags, fields, time.Now().UTC())
	bp.AddPoint(p)

	// Write the batch
	if err := db.Write(bp); err != nil {
		log.Fatal("InfluxDB write error", err)
		os.Exit(1)
	}
}

func RunQuery(db InfluxClient, command, database, precision string) (res []client.Result, err error) {
	q := client.NewQuery(command, database, precision)
	r, err := db.Query(q)
	if err != nil {
		return
	}
	if r.Error() != nil {
		err = r.Error()
		return
	}
	return r.Results, nil
}
