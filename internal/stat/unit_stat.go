package stat

import (
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type UnitStat struct {
	Name                  string
	Parallel              int
	Limit                 int
	Success               int
	Total                 int
	QPS                   float64
	Code                  map[string]int
	Elapse                time.Duration
	Rate                  float64
	ResTime               time.Duration
	StartTime             time.Time
	EndTime               time.Time
	TotalElapse           time.Duration
	IsErr                 bool
	Err                   string
	UnitStages            []*UnitStageStat
	UnitStageMilliseconds int64
	UnitStageTimes        int
	MaxStepSize           int
	Steps                 []*StepStat
	QuantileKeys          []string
	Quantile              map[string]time.Duration

	currentStage *UnitStageStat
}

func NewUnitStat(
	name string, parallel int, limit int,
	stageSeconds int64, stageTimes int, stageNumber int,
	quantile []float64, maxStepSize int,
) *UnitStat {
	if stageNumber == 0 {
		stageNumber = 100
	}

	if len(quantile) == 0 {
		quantile = []float64{80, 90, 95, 99, 99.9}
	}
	var quantileKeys []string
	for _, key := range quantile {
		quantileKeys = append(quantileKeys, strconv.FormatFloat(key, 'f', -1, 64))
	}

	if maxStepSize == 0 {
		maxStepSize = 200000
	}

	unitStageMilliseconds := stageSeconds * 1000 / int64(stageNumber)
	if unitStageMilliseconds > 0 && unitStageMilliseconds < 100 {
		unitStageMilliseconds = 100
	}

	return &UnitStat{
		Name:                  name,
		Parallel:              parallel,
		Limit:                 limit,
		UnitStageMilliseconds: unitStageMilliseconds,
		UnitStageTimes:        stageTimes,
		currentStage:          NewUnitStageStat(),
		StartTime:             time.Now(),
		MaxStepSize:           maxStepSize,
		QuantileKeys:          quantileKeys,
		Quantile:              map[string]time.Duration{},
		Code:                  map[string]int{},
	}
}

func (s *UnitStat) AddStepStat(step *StepStat) {
	s.Total += 1
	if step.Success {
		s.Success += 1
		s.Elapse += step.Elapse
	} else {
		s.Code[step.Code] += 1
	}

	s.currentStage.AddStepStat(step)
	if s.UnitStageMilliseconds != 0 && time.Now().Sub(s.currentStage.Time).Milliseconds() >= s.UnitStageMilliseconds {
		s.currentStage.Summary()
		s.UnitStages = append(s.UnitStages, s.currentStage)
		s.currentStage = NewUnitStageStat()
	}
	if s.UnitStageTimes != 0 && s.currentStage.Total >= s.UnitStageTimes {
		s.currentStage.Summary()
		s.UnitStages = append(s.UnitStages, s.currentStage)
		s.currentStage = NewUnitStageStat()
	}

	if s.MaxStepSize == 0 {
		s.Steps = append(s.Steps, step)
	} else if s.MaxStepSize > len(s.Steps) {
		s.Steps = append(s.Steps, step)
	} else {
		s.Steps[rand.Intn(len(s.Steps))] = step
	}
}

func (s *UnitStat) Summary() {
	s.EndTime = time.Now()
	s.TotalElapse = s.EndTime.Sub(s.StartTime)
	s.QPS = float64(s.Success) / s.TotalElapse.Seconds()
	if s.Success != 0 {
		s.ResTime = s.Elapse / time.Duration(s.Success)
	}
	if s.Total != 0 {
		s.Rate = float64(s.Success) / float64(s.Total)
	}
	s.Code["OK"] = s.Success

	sort.Slice(s.Steps, func(i, j int) bool {
		return s.Steps[i].Elapse < s.Steps[j].Elapse
	})

	if len(s.Steps) == 0 {
		for _, key := range s.QuantileKeys {
			s.Quantile[key] = 0
		}
	} else {
		for _, key := range s.QuantileKeys {
			v, _ := strconv.ParseFloat(key, 64)
			s.Quantile[key] = s.Steps[uint(float64(len(s.Steps))*v/100.0)].Elapse
		}
	}
}
