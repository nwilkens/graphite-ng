package main

import (
	"github.com/Dieterbe/graphite-ng/chains"
	"github.com/Dieterbe/graphite-ng/metrics"
)

func ReadMetric(name string) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	data := map[string]*metrics.Metric{
		"stats.web1.bytes_received": metrics.NewMetric(
			"stats.web1.bytes_received",
			[]*metrics.Datapoint{
				metrics.NewDatapoint(120, 3, true),
				metrics.NewDatapoint(180, 0, false),
				metrics.NewDatapoint(240, 2, true),
				metrics.NewDatapoint(300, 0, true),
				metrics.NewDatapoint(360, 1, true),
				metrics.NewDatapoint(420, 1, true),
				metrics.NewDatapoint(480, 1.5, true),
				metrics.NewDatapoint(540, 2, true),
			},
		),
		"stats.web2.bytes_received": metrics.NewMetric(
			"stats.web1.bytes_received",
			[]*metrics.Datapoint{
				metrics.NewDatapoint(120, 4, true),
				metrics.NewDatapoint(180, 1, true),
				metrics.NewDatapoint(240, 2, true),
				metrics.NewDatapoint(300, 0, false),
				metrics.NewDatapoint(360, 1, true),
				metrics.NewDatapoint(420, 1, true),
				metrics.NewDatapoint(480, 1, true),
				metrics.NewDatapoint(540, 3, true),
			},
		),
	}
	metric := data[name]
	go func(our_el chains.ChainEl, metric *metrics.Metric) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		// if we don't have enough data to cover the requested timespan, fill with nils
		if metric.Data[0].Ts > from {
			for new_ts := from; new_ts < metric.Data[0].Ts; new_ts += 60 {
				our_el.Link <- *metrics.NewDatapoint(new_ts, 0.0, false)
			}
		}
		for _, d := range metric.Data {
			if d.Ts >= from && until <= until {
				our_el.Link <- *d
			}
		}
		if metric.Data[len(metric.Data)-1].Ts < until {
			for new_ts := metric.Data[len(metric.Data)-1].Ts + 60; new_ts <= until+60; new_ts += 60 {
				our_el.Link <- *metrics.NewDatapoint(new_ts, 0.0, false)
			}
		}
	}(our_el, metric)
	return
}
