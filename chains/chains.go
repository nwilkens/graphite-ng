package chains

import (
	"github.com/graphite-ng/graphite-ng/metrics"
)

type ChainEl struct {
	Settings chan int32             // dependent fn will send the from/until it needs to dependency function
	Link     chan metrics.Datapoint // dependency fn will send values to dependent fn
}

func NewChainEl() *ChainEl {
	settings := make(chan int32)
	link := make(chan metrics.Datapoint)
	return &ChainEl{settings, link}
}
