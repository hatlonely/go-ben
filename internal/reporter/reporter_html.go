package reporter

import "github.com/hatlonely/go-ben/internal/stat"

type HtmlReporterOptions struct {
	Font struct {
		Style   string
		Body    string
		Code    string
		Echarts string
	}
	Extra struct {
		Head       string
		BodyHeader string
		BodyFooter string
	}
	Padding struct {
		X int
		Y int
	}
}

func NewHtmlReporterWithOptions(options *HtmlReporterOptions) (*HtmlReporter, error) {
	return &HtmlReporter{
		options: options,
	}, nil
}

type HtmlReporter struct {
	options *HtmlReporterOptions
}

func (r *HtmlReporter) Report(test *stat.TestStat) string {
	return "html"
}
