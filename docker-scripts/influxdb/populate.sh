#!/bin/sh

set -e

influx -execute "CREATE DATABASE $DOCKER_INFLUXDB_NAME"
