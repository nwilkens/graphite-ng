package stores

import (
	"errors"
	"fmt"
	"github.com/graphite-ng/graphite-ng/chains"
	"github.com/graphite-ng/graphite-ng/util"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"github.com/stvp/go-toml-config"
	"strconv"
)

type Es struct {
	es_host        string
	es_port        int
	es_max_pending int
	in_port        int
}

func NewEs(path string) *Es {
	var (
		es_host        = config.String("elasticsearch.host", "undefined")
		es_port        = config.Int("elasticsearch.port", 9200)
		es_max_pending = config.Int("elasticsearch.max_pending", 1000000)
		in_port        = config.Int("in.port", 2003)
	)
	err := config.Parse(path)
	util.DieIfError(err)
	api.Domain = *es_host
	api.Port = strconv.Itoa(*es_port)
	return &Es{*es_host, *es_port, *es_max_pending, *in_port}
}

func IsEsMetric(name string) (found bool, err error) {
	out, err := core.SearchUri("carbon-es", "datapoint", fmt.Sprintf("metric:%s&size=1", name), "", 0)
	if err != nil {
		return false, errors.New(fmt.Sprintf("error checking ES for %s: %s", name, err.Error()))
	}
	return (out.Hits.Total > 0), nil
}

func ReadEsMetric(name string) (our_el *chains.ChainEl, err error) {
	our_el = chains.NewChainEl()
	go func(our_el *chains.ChainEl) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		qry := map[string]interface{}{
			"query": map[string]interface{}{
				"term": map[string]string{"metric": name},
				"range": map[string]interface{}{
					"ts": map[string]string{"from": strconv.Itoa(int(from)), "to": strconv.Itoa(int(until))},
				},
			},
		}
		// { "bool": { "must": [ {"term": ... }, {"range": ...}] }}

		// TODO: sorting?
		out, err := core.SearchRequest(true, "carbon-es", "datapoint", qry, "", 0)
		if err != nil {
			panic(fmt.Sprintf("error reading ES for %s: %s", name, err.Error()))

		}
		// if we don't have enough data to cover the requested timespan, fill with nils
		/* if metric.Data[0].Ts > from {
			for new_ts := from; new_ts < metric.Data[0].Ts; new_ts += 60 {
				our_el.Link <- *metrics.NewDatapoint(new_ts, 0.0, false)
			}
		}
		for _, d := range metric.Data {
			if d.Ts >= from && until <= until {
				our_el.Link <- *d
			}
		}
		if metric.Data[len(metric.Data)-1].Ts < until {
			for new_ts := metric.Data[len(metric.Data)-1].Ts + 60; new_ts <= until+60; new_ts += 60 {
				our_el.Link <- *metrics.NewDatapoint(new_ts, 0.0, false)
			}
		}
		*/

		fmt.Println(out)
	}(our_el)
	return our_el, nil
}
