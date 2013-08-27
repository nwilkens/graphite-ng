package main

import (
	"github.com/Dieterbe/graphite-ng/chains"
	"github.com/Dieterbe/graphite-ng/metrics"
)

// like with graphite, it is assumed datapoints from different inputs are time synchronized
// at some point we might lift that and take it into account in individual functions
func FnSum(dep_els ...chains.ChainEl) (our_el chains.ChainEl) {
	our_el = *chains.NewChainEl()
	go func(our_el chains.ChainEl, dep_els []chains.ChainEl) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		for _, dep_el := range dep_els {
			dep_el.Settings <- from
			dep_el.Settings <- until
		}
		var sum float64
		var known bool
		// for every point in time (can't iterate over them here, they come from the channels)
		for {
			// sum the datapoints from the different channels together (each dp from each chan is one term)
			// we're done when we reached the last channel and the ts == until
			// if one or more of the points is !known, the resulting sum is not known
			for i, c := range dep_els {
				// first term in the sum, reset the data that will go into datapoint
				if i == 0 {
					known = true
					sum = 0.0
				}
				d := <-c.Link
				if known {
					if !d.Known {
						known = false
						our_el.Link <- *metrics.NewDatapoint(d.Ts, 0.0, false)
						if i == len(dep_els)-1 && d.Ts >= until {
							return
						}
					} else {
						sum += d.Value
						if i == len(dep_els)-1 {
							our_el.Link <- *metrics.NewDatapoint(d.Ts, sum, true)
							if d.Ts >= until {
								return
							}
						}
					}
				}
			}
		}
	}(our_el, dep_els)
	return
}

// todo: allow N inputs and outputs
func FnScale(dep_el chains.ChainEl, multiplier float64) (our_el chains.ChainEl) {
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

func FnDerivative(dep_el chains.ChainEl) (our_el chains.ChainEl) {
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

var Functions = map[string]string{
	"sum":        "FnSum",
	"sumSeries":  "FnSum",
	"scale":      "FnScale",
	"derivative": "FnDerivative",
}
