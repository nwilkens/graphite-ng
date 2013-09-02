package main

import (
	"github.com/Dieterbe/graphite-ng/chains"
	"github.com/Dieterbe/graphite-ng/stores"
)

func ReadMetric(name string) (our_el chains.ChainEl) {
	is_text := stores.IsTextMetric(name)
	if is_text {
		our_el, err := stores.ReadTextMetric(name)
		if err == nil {
			return *our_el
		} else {
			panic("Error reading metric " + name + " from TextStore: " + err.Error())
		}
	} else {
		panic("Could not find metric " + name + " in any of the stores")
	}
}
