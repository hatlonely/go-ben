package framework

type GroupDesc struct {
	Times    int
	Seconds  int
	Parallel []int
}

type UnitDesc struct {
	Name string
	Seed map[string]string
}

type PlanDesc struct {
	Name  string
	Group []GroupDesc
	Unit  []UnitDesc
}
