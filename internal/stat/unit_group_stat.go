package stat

import "strconv"

type UnitGroupStat struct {
	Idx      int
	Seconds  int
	Times    int
	Units    []*UnitStat
	Quantile []string
	Monitor  map[string]*MonitorStat
	IsErr    bool
	Err      string
}

func NewUnitGroupStat(idx int, seconds int, times int, quantile []float64) *UnitGroupStat {
	if len(quantile) == 0 {
		quantile = []float64{80, 90, 95, 99, 99.9}
	}
	var quantileKeys []string
	for _, key := range quantile {
		quantileKeys = append(quantileKeys, strconv.FormatFloat(key, 'f', -1, 64))
	}

	return &UnitGroupStat{
		Idx:      idx,
		Seconds:  seconds,
		Times:    times,
		Quantile: quantileKeys,
	}
}

func (s *UnitGroupStat) AddUnitStat(unit *UnitStat) {
	s.Units = append(s.Units, unit)
}

func (s *UnitGroupStat) AddMonitorStat(name string, monitor *MonitorStat) {
	s.Monitor[name] = monitor
}
