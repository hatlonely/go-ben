package stat

import "fmt"

type PlanStat struct {
	ID          string
	Description string
	Command     string
	Name        string
	IsErr       bool
	Err         string
	UnitGroups  []*UnitGroupStat
}

func NewPlanStat(id string, name string) *PlanStat {
	return &PlanStat{
		ID:   id,
		Name: name,
	}
}

func (s *PlanStat) SetError(err error) *PlanStat {
	s.IsErr = true
	s.Err = fmt.Sprintf("%+v", err)
	return s
}

func (s *PlanStat) AddUnitGroupStat(unitGroup *UnitGroupStat) {
	s.UnitGroups = append(s.UnitGroups, unitGroup)
}
