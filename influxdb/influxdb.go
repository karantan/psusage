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

// AddPoint creates and writes a Influx point to the `db` InfluxClient (InfluxDB).
// It will retry max of 3 times in case of an error.
func AddPoint(db InfluxClient, u collect.CPU_Usage, hostname string) {
	point, err := CreatePoint(u, hostname)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for retry := 0; retry < 3; retry++ {
		err := WritePoint(db, point, retry)
		if err == nil {
			break
		}

		log.Error("InfluxDB write error", err)
		time.Sleep(time.Duration(retry) * time.Second)
	}
	log.Error("Max retries exceeded while trying to write to InfluxDB")
	os.Exit(1)
}

// CreatePoint creates a `collect.CPU_Usage` point InfluxDB point which we can send
// to the InfluxDB (see WritePoint).
func CreatePoint(u collect.CPU_Usage, hostname string) (*client.Point, error) {
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
	return client.NewPoint(measurement, tags, fields, time.Now().UTC())
}

// WritePoint writes influx points in batch.
func WritePoint(db InfluxClient, p *client.Point, retry int) error {
	// Make sure the database exists (`CREATE DATABASE psusage`)
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "psusage",
		Precision: "s",
	})
	bp.AddPoint(p)

	// Write the batch
	return db.Write(bp)
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
