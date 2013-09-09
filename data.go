package main

import (
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/stores"
)

func ReadMetric(name string) (our_el chains.ChainEl) {
	var found bool
	var err error
	found, err = stores.IsTextMetric(name)
	if err != nil {
		panic("Error reading metric " + name + " from TextStore: " + err.Error())
	}
	if found {
		our_el, err := stores.ReadTextMetric(name)
		if err == nil {
			return *our_el
		} else {
			panic("Error reading metric " + name + " from TextStore: " + err.Error())
		}
	}

	found, err = stores.IsEsMetric(name)
	if err != nil {
		panic("Error reading metric " + name + " from EsStore: " + err.Error())
	}
	if found {
		our_el, err := stores.ReadEsMetric(name)
		if err == nil {
			return *our_el
		} else {
			panic("Error reading metric " + name + " from EsStore: " + err.Error())
		}
	}
	panic("Could not find metric " + name + " in any of the stores")
}
