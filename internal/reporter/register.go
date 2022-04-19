package reporter

func init() {
	RegisterReporter("Json", NewJsonReporterWithOptions)
	RegisterReporter("json", NewJsonReporterWithOptions)
	RegisterReporter("Html", NewHtmlReporterWithOptions)
	RegisterReporter("html", NewHtmlReporterWithOptions)
	RegisterReporter("Text", NewTextReporterWithOptions)
	RegisterReporter("text", NewTextReporterWithOptions)
}
