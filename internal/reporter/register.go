package reporter

func init() {
	RegisterReporter("Json", NewJsonReporterWithOptions)
	RegisterReporter("Html", NewHtmlReporterWithOptions)
	RegisterReporter("Text", NewTextReporterWithOptions)
}
