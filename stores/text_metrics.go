package stores

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Dieterbe/graphite-ng/chains"
	"github.com/Dieterbe/graphite-ng/metrics"
	"os"
	"strconv"
	"strings"
)

func GetTextMetricPath(name string) string {
	return fmt.Sprintf("text_metrics/%s.txt", name)
}

func IsTextMetric(name string) bool {
	_, err := os.Stat(GetTextMetricPath(name))
	return (err == nil)
}

func ReadTextMetric(name string) (our_el *chains.ChainEl, err error) {
	var file *os.File
	path := GetTextMetricPath(name)
	if file, err = os.Open(path); err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	datapoints := make([]*metrics.Datapoint, 0)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		ts, _ := strconv.ParseInt(parts[0], 10, 32)
		val, _ := strconv.ParseFloat(parts[1], 64)
		known, _ := strconv.ParseBool(parts[2])
		dp := metrics.NewDatapoint(int32(ts), val, known)
		datapoints = append(datapoints, dp)
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error reading %s: %s", path, err))
	}
	metric := metrics.NewMetric(name, datapoints)

	our_el = chains.NewChainEl()
	go func(our_el *chains.ChainEl, metric *metrics.Metric) {
		from := <-our_el.Settings
		until := <-our_el.Settings
		// if we don't have enough data to cover the requested timespan, fill with nils
		if metric.Data[0].Ts > from {
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
	}(our_el, metric)
	return our_el, nil
}
