package framework

type Options struct {
	TestDirectory string
	PlanDirectory string
	Customize     string
	Reporter      string
	X             string
	JsonStat      string
	Hook          string
	Lang          string
}

func NewFrameworkWithOptions(options *Options) (*Framework, error) {
	return &Framework{
		options: options,
	}, nil
}

type Framework struct {
	options *Options
}
