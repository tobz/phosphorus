package statistics

import "fmt"
import "time"
import "strconv"
import "github.com/rcrowley/go-metrics"
import "github.com/rcrowley/go-metrics/influxdb"
import "github.com/tobz/phosphorus/interfaces"

var Registry metrics.Registry = metrics.NewRegistry()

func init() {
	// Capture runtime statistics every 5 seconds.
	metrics.RegisterRuntimeMemStats(Registry)
	go metrics.CaptureRuntimeMemStats(Registry, time.Second*5)
}

func ConfigureInfluxDB(c interfaces.Config) error {
	// Grab all the relevant configuration bits.
	influxHost, err := c.GetAsString("statistics/host")
	if err != nil {
		return fmt.Errorf("caught an error trying to read InfluxDB host: %s", err)
	}

	influxUsername, err := c.GetAsString("statistics/username")
	if err != nil {
		return fmt.Errorf("caught an error trying to read InfluxDB username: %s", err)
	}

	influxPassword, err := c.GetAsString("statistics/password")
	if err != nil {
		return fmt.Errorf("caught an error trying to read InfluxDB password: %s", err)
	}

	influxDatabase, err := c.GetAsString("statistics/database")
	if err != nil {
		return fmt.Errorf("caught an error trying to read InfluxDB database: %s", err)
	}

	flushRateRaw, err := c.GetAsString("statistics/flushRate")
	if err != nil {
		return fmt.Errorf("caught an error trying to read metrics flush rate: %s", err)
	}

	flushRateFloat, err := strconv.ParseFloat(flushRateRaw, 64)
	if err != nil {
		return fmt.Errorf("caught an error trying to parse flush rate: %s", err)
	}

	flushRate := int64(flushRateFloat)

	// We should be good to go.
	go influxdb.Influxdb(Registry, time.Duration(flushRate), &influxdb.Config{
		Host:     influxHost,
		Database: influxDatabase,
		Username: influxUsername,
		Password: influxPassword,
	})

	return nil
}
