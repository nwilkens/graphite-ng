package metrics

import (
	"fmt"
)

type Datapoint struct {
	Ts    int32
	Value float64
	Known bool
}
type Metric struct {
	Name string
	Data []*Datapoint
}

func NewDatapoint(ts int32, value float64, known bool) *Datapoint {
	return &Datapoint{ts, value, known}
}

func (d *Datapoint) String() string {
	return fmt.Sprintf("datapoint ts=%d, value=%f, known=%t", d.Ts, d.Value, d.Known)
}

func NewMetric(name string, data []*Datapoint) *Metric {
	return &Metric{name, data}
}
