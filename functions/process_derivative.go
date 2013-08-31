package functions

import (
	"github.com/Dieterbe/graphite-ng/chains"
	"github.com/Dieterbe/graphite-ng/metrics"
)

func init() {
	Functions["derivative"] = "ProcessDerivative"
}
func ProcessDerivative(dep_el chains.ChainEl) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	go func(our_el chains.ChainEl, dep_el chains.ChainEl) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		dep_el.Settings <- from - 60
		dep_el.Settings <- until
		var last_dp *metrics.Datapoint

		for {
			d := <-dep_el.Link
			if last_dp == nil {
				last_dp = &d
				continue
			}
			if d.Known && last_dp.Known {
				our_el.Link <- *metrics.NewDatapoint(d.Ts, d.Value-last_dp.Value, true)
			} else {
				our_el.Link <- *metrics.NewDatapoint(d.Ts, 0.0, false)
			}
			last_dp = &d
			if d.Ts >= until {
				return
			}
		}
	}(our_el, dep_el)
	return
}
