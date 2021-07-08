package sabnzbd_prometheus

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// SabnzbdCollector is a custom collector.
// It handles retrieving the stats from sabnzbd and creating the associated prometheus metrics
// The reason to have a collector is for 2 reasons:
// (1) Collect metrics at /metrics query time
// (2) Wanted to provide the `server` as a label
// the builtIn `prom.CounterFunc` (and other `nFunc` metrics) will have to query the api many times to achive this
type SabNzbdCollector struct {
	client                    SabNzbdClient
	downloadedOverallMetric   *prometheus.Desc
	downloadedServerMetric    *prometheus.Desc
	queueSizeMetric           *prometheus.Desc
	queueBytesPerSecondMetric *prometheus.Desc
	queueRemainingBytes       *prometheus.Desc
	queueTotalSizeBytes       *prometheus.Desc
	queueRemainingTimeSeconds *prometheus.Desc
	speedLimitPercentage      *prometheus.Desc
	speedLimitAbsolute        *prometheus.Desc
}

func NewSabNzbdCollector(baseUri string, apiKey string) *SabNzbdCollector {
	client := NewSabNzbdClient(baseUri, apiKey)
	return &SabNzbdCollector{
		client: client,
		downloadedOverallMetric: prometheus.NewDesc(
			"total_downloaded",
			"SabNzbd Overall total number of bytes downloaded",
			[]string{"period"},
			prometheus.Labels{},
		),
		downloadedServerMetric: prometheus.NewDesc(
			"server_total_downloaded",
			"SabNzbd per server total number of bytes downloaded",
			[]string{"period", "server"},
			prometheus.Labels{},
		),
		queueSizeMetric: prometheus.NewDesc(
			"queue_size",
			"Remaining number of items in queue",
			nil,
			prometheus.Labels{},
		),
		queueBytesPerSecondMetric: prometheus.NewDesc(
			"queue_download_bytes_per_second",
			"Current download rate in bytes/second",
			nil,
			prometheus.Labels{},
		),
		queueRemainingBytes: prometheus.NewDesc(
			"queue_remaining_bytes",
			"Remaining number of bytes in queue",
			nil,
			prometheus.Labels{},
		),
		queueTotalSizeBytes: prometheus.NewDesc(
			"queue_total_size_bytes",
			"Total number of bytes in queue",
			nil,
			prometheus.Labels{},
		),
		queueRemainingTimeSeconds: prometheus.NewDesc(
			"queue_remaining_time_seconds",
			"Remaining time in seconds in queue",
			nil,
			prometheus.Labels{},
		),
		speedLimitPercentage: prometheus.NewDesc(
			"speed_limit_used_percentage",
			"Percentage of speed limit used",
			nil,
			prometheus.Labels{},
		),
		speedLimitAbsolute: prometheus.NewDesc(
			"speed_limit_absolute",
			"Speed limit in bytes",
			nil,
			prometheus.Labels{},
		),
	}

}
func (s *SabNzbdCollector) Describe(c chan<- *prometheus.Desc) {
	fmt.Println("descriving")
	c <- s.downloadedOverallMetric
	c <- s.downloadedServerMetric
	c <- s.queueSizeMetric
	c <- s.queueBytesPerSecondMetric
	c <- s.queueRemainingBytes
	c <- s.queueTotalSizeBytes
	c <- s.queueRemainingTimeSeconds
	c <- s.speedLimitPercentage
	c <- s.speedLimitAbsolute
}

func (s *SabNzbdCollector) Collect(c chan<- prometheus.Metric) {
	fmt.Println("collecting")
	statsResponse, err := s.client.GetServerStats()
	if err != nil {
		fmt.Printf("error retrieving SabNzbd server stats: %e\n", err)
		return
	}
	// Overall server metrics
	c <- prometheus.MustNewConstMetric(s.downloadedOverallMetric, prometheus.GaugeValue, float64(statsResponse.Total), "total")
	c <- prometheus.MustNewConstMetric(s.downloadedOverallMetric, prometheus.GaugeValue, float64(statsResponse.Month), "month")
	c <- prometheus.MustNewConstMetric(s.downloadedOverallMetric, prometheus.GaugeValue, float64(statsResponse.Week), "week")
	c <- prometheus.MustNewConstMetric(s.downloadedOverallMetric, prometheus.GaugeValue, float64(statsResponse.Day), "day")

	// per server metrics
	for server, metric := range statsResponse.Servers {
		// todo: dont use must and error
		c <- prometheus.MustNewConstMetric(s.downloadedServerMetric, prometheus.GaugeValue, float64(metric.Total), "total", server)
		c <- prometheus.MustNewConstMetric(s.downloadedServerMetric, prometheus.GaugeValue, float64(metric.Month), "month", server)
		c <- prometheus.MustNewConstMetric(s.downloadedServerMetric, prometheus.GaugeValue, float64(metric.Week), "week", server)
		c <- prometheus.MustNewConstMetric(s.downloadedServerMetric, prometheus.GaugeValue, float64(metric.Day), "day", server)
	}

	queueResponse, err := s.client.GetQueue()
	if err != nil {
		fmt.Printf("error retrieving SabNzbd queue stats: %e\n", err)
		return
	}
	queue := queueResponse.Queue

	c <- prometheus.MustNewConstMetric(s.queueSizeMetric, prometheus.GaugeValue, float64(queue.NoOfSlotsTotal))

	bps, err := strconv.ParseFloat(queue.KbPerSec, 64)
	if err != nil {
		printParseFloatError("kbpersec", queue.KbPerSec, err)
		return
	}
	c <- prometheus.MustNewConstMetric(s.queueBytesPerSecondMetric, prometheus.GaugeValue, float64(bps*1024))

	qrs, err := strconv.ParseFloat(queue.MbLeft, 64)
	if err != nil {
		printParseFloatError("mbleft", queue.MbLeft, err)
		return
	}
	c <- prometheus.MustNewConstMetric(s.queueRemainingBytes, prometheus.GaugeValue, float64(qrs*1024*1024))

	qts, err := strconv.ParseFloat(queue.Mb, 64)
	if err != nil {
		printParseFloatError("mb", queue.Mb, err)
		return
	}
	c <- prometheus.MustNewConstMetric(s.queueTotalSizeBytes, prometheus.GaugeValue, float64(qts*1024*1024))

	tr, err := time.Parse("15:04:05", queue.TimeLeft)
	if err != nil {
		printParseFloatError("timeleft", queue.TimeLeft, err)
		return
	}
	// TODO: should probably just use a duration
	timeRemaining := (tr.Hour() * 3600) + (tr.Minute() * 60) + tr.Second()
	c <- prometheus.MustNewConstMetric(s.queueRemainingTimeSeconds, prometheus.GaugeValue, float64(timeRemaining))

}

func printParseFloatError(key string, value string, err error) {
	fmt.Printf("error converting `%s` value %s to time; %e\n", key, value, err)
}
