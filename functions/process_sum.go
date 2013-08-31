package functions

import (
	"github.com/Dieterbe/graphite-ng/chains"
	"github.com/Dieterbe/graphite-ng/metrics"
)

func init() {
	Functions["sum"] = "ProcessSum"
	Functions["sumSeries"] = "ProcessSum"
}

// like with graphite, it is assumed datapoints from different inputs are time synchronized
// at some point we might lift that and take it into account in individual functions
func ProcessSum(dep_els ...chains.ChainEl) (our_el chains.ChainEl) {
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
