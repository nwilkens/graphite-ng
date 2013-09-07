package es

type Datapoint struct {
	Metric string  `json:"metric"`
	Ts     int32   `json:"ts"`
	Value  float64 `json:"value"`
}
