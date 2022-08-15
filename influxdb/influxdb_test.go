//go:build integration
// +build integration

package influxdb

import (
	"encoding/json"
	"os"
	"psusage/collect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getInfluxDSN() string {
	influxDSN := os.Getenv("DEV_INFLUX_DSN") // see Dockerfile for this env
	if influxDSN == "" {
		influxDSN = "http://127.0.0.1:8086"
	}
	return influxDSN
}

func TestNewInfluxDSN(t *testing.T) {
	influx := NewInfluxDSN(getInfluxDSN())
	_, _, err := influx.Ping(5)
	assert.NoError(t, err)
}

func TestAddPoint(t *testing.T) {
	influx := NewInfluxDSN(getInfluxDSN())
	p := collect.CPU_Usage{
		PCPU:     0.1,
		Duration: 5,
		Program:  "myprogram",
		PID:      100,
		User:     "worker",
	}
	AddPoint(influx, p, "server")

	wantMeasurements := "cpu"
	wantColumns := []string{"time", "duration", "pcpu", "program", "server", "user"}
	wantValues := []any{json.Number("5"), json.Number("0.1"), "myprogram", "server", "worker"}

	got, err := RunQuery(influx, "SELECT * FROM cpu", "psusage", "s")
	assert.Equal(t, wantMeasurements, got[0].Series[0].Name)
	assert.Equal(t, wantColumns, got[0].Series[0].Columns)
	assert.Equal(t, wantValues, got[0].Series[0].Values[0][1:6])
	assert.NoError(t, err)
}
