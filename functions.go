package main

import (
	"fmt"
)

// like with graphite, it is assumed datapoints from different inputs are time synchronized
// at some point we might lift that and take it into account in individual functions
func FnSum(from int32, until int32, in ...chan Datapoint) chan Datapoint {
	out := make(chan Datapoint)
	go func(from int32, until int32, out chan Datapoint, in []chan Datapoint) {
		var sum float64
		// for every point in time (can't iterate over them here, they come from the channels)
		for {
			// sum the datapoints from the different channels together (each dp from each chan is one term)
			// we're done when we reached the last channel and the ts == until
			for i, c := range in {
				// first term in the sum, reset sum
				if i == 0 {
					//fmt.Println("reset sum")
					sum = 0.0
				}
				d := <-c
				fmt.Println("FnSum read", d, "from chan", i)
				// if one of the values in the sum is unknown, the result is unknown
				if !d.known {
					//fmt.Println("FnSum yielding nil")
					out <- *NewDatapoint(d.ts, 0.0, false)
					if i == len(in)-1 && d.ts == until {
						return
					}
					break
				}
				sum += d.value
				fmt.Println("sum is now", sum)
				if i == len(in)-1 {
					//fmt.Println("FnSum yielding a sum", sum)
					out <- *NewDatapoint(d.ts, sum, true)
					if d.ts == until {
						return
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
			fmt.Println("FnScale read", d)
			if !d.known {
				out <- *NewDatapoint(d.ts, 0.0, false)
				if d.ts == until {
					return
				}
				continue
			}
			//fmt.Println("FnScale yielding...")
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
