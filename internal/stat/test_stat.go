package stat

type TestStat struct {
	ID          string
	Directory   string
	Name        string
	Description string
	IsErr       bool
	Err         string
	Plans       []*PlanStat
	SubTests    []*TestStat
}

func NewTestStat(id string, directory string, name string, description string) *TestStat {
	return &TestStat{
		ID:          id,
		Directory:   directory,
		Name:        name,
		Description: description,
	}
}

func (s *TestStat) SetError(err error) *TestStat {
	s.IsErr = true
	s.Err = err.Error()
	return s
}

func (s *TestStat) AddPlanStat(plan *PlanStat) {
	s.Plans = append(s.Plans, plan)
}

func (s *TestStat) AddSubTestStat(subTest *TestStat) {
	s.SubTests = append(s.SubTests, subTest)
}
