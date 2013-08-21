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
				NewDatapoint(0, 1.5, true),
				NewDatapoint(60, 2, true),
				NewDatapoint(120, 3, true),
				NewDatapoint(180, 0, false),
				NewDatapoint(240, 2, true),
				NewDatapoint(300, 0, true),
				NewDatapoint(360, 1, true),
				NewDatapoint(420, 1, true),
			},
		),
		"stats.web2.bytes_received": NewMetric_data(
			"stats.web1.bytes_received",
			[]*Datapoint{
				NewDatapoint(0, 1, true),
				NewDatapoint(60, 3, true),
				NewDatapoint(120, 4, true),
				NewDatapoint(180, 1, true),
				NewDatapoint(240, 2, true),
				NewDatapoint(300, 0, false),
				NewDatapoint(360, 1, true),
				NewDatapoint(420, 1, true),
			},
		),
	}
	metric := data[name]
	out := make(chan Datapoint)
	go func(out chan Datapoint, metric *Metric_data) {
		for _, d := range metric.data {
			if d.ts >= from && until <= until {
				fmt.Printf("ReadMetric %s writing %s\n", name, *d)
				out <- *d
			}
		}
	}(out, metric)
	return out

}
