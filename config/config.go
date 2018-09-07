package config

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

const (
	Perc50  = float64(0.50)
	Perc75  = float64(0.50)
	Perc95  = float64(0.50)
	Perc99  = float64(0.99)
	Perc999 = float64(0.999)
	Perc100 = float64(1)
)

type PutMetricsClient interface {
	PutMetricData(*cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error)
}

type Config struct {
	Filter                Filter
	Client                PutMetricsClient
	ReportingInterval     time.Duration
	Namespace             string
	StaticDimensions      map[string]string
	ResetCountersOnReport bool
	Silence               bool
}

type Filter interface {
	ShouldReport(metric string, value float64, config Config) bool
	Percentiles(metric string) []float64
}

type NoFilter struct{}

func (n *NoFilter) ShouldReport(metric string, value float64, cfg Config) bool {
	if !cfg.Silence {
		log.Printf("at=should-report metric=%s ", metric)
	}
	return true
}

func (n *NoFilter) Percentiles(metric string) []float64 {
	return []float64{Perc50, Perc75, Perc95, Perc99, Perc999, Perc100}
}

type AllFilter struct{}

func (n *AllFilter) ShouldReport(metric string, value float64, cfg Config) bool {
	if !cfg.Silence {
		log.Printf("at=no-report metric=%s ", metric)
	}
	return false
}

func (n *AllFilter) Percentiles(metric string) []float64 {
	return []float64{}
}

/*
type DynamoDBFilter struct {
	globalEnabledMetrics []string
	perInstanceEnabledMetrics map[string]string
}

func (d *DynamodbConfig) PollConfig() {
	poll once every few minutes, read enabled metrics
}
*/
