package main

// like with graphite, it is assumed datapoints from different inputs are time synchronized
// at some point we might lift that and take it into account in individual functions
func FnSum(from int32, until int32, in ...chan Datapoint) chan Datapoint {
	out := make(chan Datapoint)
	go func(from int32, until int32, out chan Datapoint, in []chan Datapoint) {
		var sum float64
		var known bool
		// for every point in time (can't iterate over them here, they come from the channels)
		for {
			// sum the datapoints from the different channels together (each dp from each chan is one term)
			// we're done when we reached the last channel and the ts == until
			// if one or more of the points is !known, the resulting sum is not known
			for i, c := range in {
				// first term in the sum, reset the data that will go into datapoint
				if i == 0 {
					known = true
					sum = 0.0
				}
				d := <-c
				if known {
					if !d.known {
						known = false
						out <- *NewDatapoint(d.ts, 0.0, false)
						if i == len(in)-1 && d.ts == until {
							return
						}
					} else {
						sum += d.value
						if i == len(in)-1 {
							out <- *NewDatapoint(d.ts, sum, true)
							if d.ts == until {
								return
							}
						}
					}
				}
			}
		}
	}(from, until, out, in)
	return out
}

// todo: allow N inputs and outputs
func FnScale(from int32, until int32, in chan Datapoint, multiplier float64) chan Datapoint {
	out := make(chan Datapoint)
	go func(from int32, until int32, out chan Datapoint, in chan Datapoint, multiplier float64) {
		for {
			d := <-in
			if !d.known {
				out <- *NewDatapoint(d.ts, 0.0, false)
				if d.ts == until {
					return
				}
				continue
			}
			out <- *NewDatapoint(d.ts, d.value*multiplier, true)
			if d.ts == until {
				return
			}
		}
	}(from, until, out, in, multiplier)
	return out
}

var Functions = map[string]string{
	"sum":       "FnSum(from, until, ",
	"sumSeries": "FnSum(from, until, ",
	"scale":     "FnScale(from, until, ",
}
