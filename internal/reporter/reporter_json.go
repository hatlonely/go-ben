package reporter

import (
	"github.com/hatlonely/go-ben/internal/stat"
	"github.com/hatlonely/go-kit/strx"
)

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
	return strx.JsonMarshalIndent(test)
}
