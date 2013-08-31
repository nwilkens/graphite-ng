package main

import (
	"fmt"
	"github.com/Dieterbe/graphite-ng/chains"
	"github.com/Dieterbe/graphite-ng/metrics"
)

func ReadMetric(name string) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	data := map[string]*metrics.Metric{
		"stats.web1": metrics.NewMetric(
			"stats.web1",
			[]*metrics.Datapoint{
				metrics.NewDatapoint(0, 0, true),
				metrics.NewDatapoint(60, 1, true),
				metrics.NewDatapoint(120, 2, true),
				metrics.NewDatapoint(180, 3, true),
				metrics.NewDatapoint(240, 4, true),
				metrics.NewDatapoint(300, 5, true),
				metrics.NewDatapoint(360, 5.5, true),
				metrics.NewDatapoint(420, 6, true),
				metrics.NewDatapoint(480, 6.5, true),
				metrics.NewDatapoint(540, 7, true),
			},
		),
		"stats.web2": metrics.NewMetric(
			"stats.web2",
			[]*metrics.Datapoint{
				metrics.NewDatapoint(120, 4, true),
				metrics.NewDatapoint(180, 3, true),
				metrics.NewDatapoint(240, 2, true),
				metrics.NewDatapoint(300, 0, false),
				metrics.NewDatapoint(360, 1, true),
				metrics.NewDatapoint(420, 1, true),
				metrics.NewDatapoint(480, 1, true),
				metrics.NewDatapoint(540, 2, true),
			},
		),
	}
	metric, ok := data[name]
	if !ok {
		panic(fmt.Sprintf("No such metric available: %s", name))
	}
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
