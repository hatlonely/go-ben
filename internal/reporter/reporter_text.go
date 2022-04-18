package reporter

import "github.com/hatlonely/go-ben/internal/stat"

type TextReporterOptions struct{}

func NewTextReporterWithOptions(options *TextReporterOptions) (*TextReporter, error) {
	return &TextReporter{
		options: options,
	}, nil
}

type TextReporter struct {
	options *TextReporterOptions
}

func (r *TextReporter) Report(test *stat.TestStat) string {
	return "text"
}
