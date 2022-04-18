package reporter

import "github.com/hatlonely/go-ben/internal/stat"

type JsonReporterOptions struct{}

func NewJsonReporterWithOptions(options *JsonReporterOptions) (*JsonReporter, error) {
	return &JsonReporter{
		options: options,
	}, nil
}

type JsonReporter struct {
	options *JsonReporterOptions
}

func (r *JsonReporter) Report(test *stat.TestStat) string {
	return "Json"
}
