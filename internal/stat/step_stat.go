package stat

import (
	"fmt"
	"time"
)

type StepStat struct {
	SubSteps []*SubStepStat
	Code     string
	Success  bool
	Elapse   time.Duration
	IsErr    bool
	Err      string
}

func NewStepStat() *StepStat {
	return &StepStat{
		Success: true,
	}
}

func (s *StepStat) SetError(err error) {
	s.IsErr = true
	s.Err = err.Error()
}

func (s *StepStat) AddSubStepStat(subStep *SubStepStat) {
	s.SubSteps = append(s.SubSteps, subStep)
	s.Elapse += subStep.Elapse
	if !subStep.Success {
		s.Success = false
		s.Code = fmt.Sprintf("%s.%s", subStep.Name, subStep.Code)
	}
}

func (s *StepStat) AddErrStat(name string, err error) {
	s.IsErr = true
	s.Err = err.Error()
	s.Success = false
	s.Code = fmt.Sprintf("%s.%s", name, "ERROR")
}
