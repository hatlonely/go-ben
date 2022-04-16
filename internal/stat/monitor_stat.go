package stat

import "time"

type Measurement struct {
	Time  time.Time
	Value float64
}

type MonitorStat struct {
	// 计量单位
	Unit map[string]string
	Stat map[string][]Measurement
}

func NewMonitorStat() *MonitorStat {
	return &MonitorStat{
		Unit: map[string]string{},
		Stat: map[string][]Measurement{},
	}
}
