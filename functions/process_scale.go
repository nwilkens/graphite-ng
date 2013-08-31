package functions

import (
	"github.com/Dieterbe/graphite-ng/chains"
	"github.com/Dieterbe/graphite-ng/metrics"
)

func init() {
	Functions["scale"] = "ProcessScale"
}

// todo: allow N inputs and outputs
func ProcessScale(dep_el chains.ChainEl, multiplier float64) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	go func(our_el chains.ChainEl, dep_el chains.ChainEl, multiplier float64) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		dep_el.Settings <- from
		dep_el.Settings <- until
		for {
			d := <-dep_el.Link
			if !d.Known {
				our_el.Link <- *metrics.NewDatapoint(d.Ts, 0.0, false)
				if d.Ts >= until {
					return
				}
				continue
			}
			our_el.Link <- *metrics.NewDatapoint(d.Ts, d.Value*multiplier, true)
			if d.Ts >= until {
				return
			}
		}
	}(our_el, dep_el, multiplier)
	return
}
