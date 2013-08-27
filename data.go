package main

import (
	"fmt"
)

type Datapoint struct {
	ts    int32
	value float64
	known bool
}
type Metric_data struct {
	name string
	data []*Datapoint
}

func NewDatapoint(ts int32, value float64, known bool) *Datapoint {
	return &Datapoint{ts, value, known}
}

func (d *Datapoint) String() string {
	return fmt.Sprintf("datapoint ts=%i, value=%f, known=%b", d.ts, d.value, d.known)
}

func NewMetric_data(name string, data []*Datapoint) *Metric_data {
	return &Metric_data{name, data}
}
func ReadMetric(name string, from int32, until int32) chan Datapoint {
	data := map[string]*Metric_data{
		"stats.web1.bytes_received": NewMetric_data(
			"stats.web1.bytes_received",
			[]*Datapoint{
				NewDatapoint(120, 3, true),
				NewDatapoint(180, 0, false),
				NewDatapoint(240, 2, true),
				NewDatapoint(300, 0, true),
				NewDatapoint(360, 1, true),
				NewDatapoint(420, 1, true),
				NewDatapoint(480, 1.5, true),
				NewDatapoint(540, 2, true),
			},
		),
		"stats.web2.bytes_received": NewMetric_data(
			"stats.web1.bytes_received",
			[]*Datapoint{
				NewDatapoint(120, 4, true),
				NewDatapoint(180, 1, true),
				NewDatapoint(240, 2, true),
				NewDatapoint(300, 0, false),
				NewDatapoint(360, 1, true),
				NewDatapoint(420, 1, true),
				NewDatapoint(480, 1, true),
				NewDatapoint(540, 3, true),
			},
		),
	}
	metric := data[name]
	out := make(chan Datapoint)
	go func(out chan Datapoint, metric *Metric_data) {
		// if we don't have enough data to cover the requested timespan, fill with nils
		if metric.data[0].ts > from {
			for new_ts := from; new_ts < metric.data[0].ts; new_ts += 60 {
				out <- *NewDatapoint(new_ts, 0.0, false)
			}
		}
		for _, d := range metric.data {
			if d.ts >= from && until <= until {
				out <- *d
			}
		}
		if metric.data[len(metric.data)-1].ts < until {
			for new_ts := metric.data[len(metric.data)-1].ts + 60; new_ts <= until+60; new_ts += 60 {
				out <- *NewDatapoint(new_ts, 0.0, false)
			}
		}
	}(out, metric)
	return out

}
