package stat

import "time"

type UnitStageStat struct {
	Time    time.Time
	Success int
	Total   int
	QPS     float64
	Rate    float64
	ResTime time.Duration
	Elapse  time.Duration
}

func NewUnitStageStat() *UnitStageStat {
	return &UnitStageStat{
		Time: time.Now(),
	}
}

func (s *UnitStageStat) AddStepStat(step *StepStat) {
	s.Total += 1
	if step.Success {
		s.Success += 1
		s.Elapse += step.Elapse
	}
}

func (s *UnitStageStat) Summary() {
	totalElapse := time.Now().Sub(s.Time)
	s.QPS = float64(s.Success) / totalElapse.Seconds()
	if s.Success != 0 {
		s.ResTime = s.Elapse / time.Duration(s.Success)
	}
	if s.Total != 0 {
		s.Rate = float64(s.Success) / float64(s.Total)
	}
}
