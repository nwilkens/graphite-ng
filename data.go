package main

type Datapoint struct {
	ts    int32
	value float64
	known bool
}
type Metric_data struct {
	name string
	data []Datapoint
}

func NewDatapoint(ts int32, value float64, known bool) *Datapoint {
	return &Datapoint{ts, value, known}
}
func NewMetric_data(name string, data []Datapoint) *Metric_data {
	return &Metric_data{name, data}
}
func ReadMetric(name string, from int32, until int32) chan Datapoint {
	data := map[string]Metric_data{
		"stats.web1.bytes_received": NewMetric_data(
			"stats.web1.bytes_received",
			[]Datapoint{
				Datapoint(0, 1.5, true),
				Datapoint(60, 2, true),
				Datapoint(120, 3, true),
				Datapoint(180, 0, false),
				Datapoint(240, 2, true),
				Datapoint(300, 0, true),
				Datapoint(360, 1, true),
				Datapoint(420, 1, true),
			},
		),
		"stats.web2.bytes_received": NewMetric_data(
			"stats.web1.bytes_received",
			[]Datapoint{
				Datapoint(0, 1, true),
				Datapoint(60, 3, true),
				Datapoint(120, 4, true),
				Datapoint(180, 1, true),
				Datapoint(240, 2, true),
				Datapoint(300, 0, false),
				Datapoint(360, 1, true),
				Datapoint(420, 1, true),
			},
		),
	}
	metric := &data[name]
	out := make(chan Datapoint)
	for _, d := range metric.data {
		if d.ts >= from && until <= until {
			out <- d
		}
	}

}
