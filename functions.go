package main

// like with graphite, it is assumed datapoints from different inputs are time synchronized
// at some point we might lift that and take it into account in individual functions
func sumSeries(from int32, until int32, in ...chan Datapoint) chan Datapoint {
	out := make(chan Datapoint)
	// for every point in time (can't iterate over them here, they come from the channels)
	for {
		// sum the datapoints from the different channels together (each dp from each chan is one term)
		// we're done when we reached the last channel and the ts == until
		for i, c := range in {
			var sum float64
			// first term in the sum, reset sum
			if i == 0 {
				sum = 0.0
			}
			d := <-c
			// if one of the values in the sum is unknown, the result is unknown
			if !d.known {
				out <- *NewDatapoint(d.ts, 0.0, false)
				if i == len(in)-1 && d.ts == until {
					return out
				}
				break
			}
			sum += d.value
			if i == len(in)-1 {
				out <- *NewDatapoint(d.ts, sum, true)
				if d.ts == until {
					return out
				}
			}
		}
	}
}

// todo: allow N inputs and outputs
func scale(from int32, until int32, in chan Datapoint, multiplier float64) chan Datapoint {
	out := make(chan Datapoint)
	for {
		d := <-in
		if !d.known {
			out <- *NewDatapoint(d.ts, 0.0, false)
			if d.ts == until {
				return out
			}
			continue
		}
		out <- *NewDatapoint(d.ts, d.value*multiplier, true)
		if d.ts == until {
			return out
		}
	}
}

var Functions = map[string]interface{}{
	"sum":       sumSeries,
	"sumSeries": sumSeries,
	"scale":     scale,
}
