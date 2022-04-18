package reporter

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/hatlonely/go-ben/internal/i18n"
	"github.com/hatlonely/go-ben/internal/stat"
	"github.com/pkg/errors"
)

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
	Lang string
	I18n i18n.I18n
}

func NewHtmlReporterWithOptions(options *HtmlReporterOptions) (*HtmlReporter, error) {
	i18n_ := i18n.NewI18n(options.Lang, &options.I18n)

	reporter := &HtmlReporter{
		options: options,
		i18n:    i18n_,
	}

	funcs := template.FuncMap{
		"RenderTest":      reporter.RenderTest,
		"RenderPlan":      reporter.RenderPlan,
		"RenderUnitGroup": reporter.RenderUnitGroup,
		"Markdown": func(text string) string {
			extensions := parser.CommonExtensions | parser.AutoHeadingIDs
			parser := parser.NewWithExtensions(extensions)
			return string(markdown.ToHTML([]byte(text), parser, nil))
		},
	}

	reporter.reportTpl = template.Must(template.New("").Funcs(funcs).Parse(reportTplStr))
	reporter.testTpl = template.Must(template.New("").Funcs(funcs).Parse(testTplStr))
	reporter.planTpl = template.Must(template.New("").Funcs(funcs).Parse(planTplStr))
	reporter.unitGroupTpl = template.Must(template.New("").Funcs(funcs).Parse(unitGroupTplStr))

	return reporter, nil
}

var reportTplStr = `<!DOCTYPE html>
<html lang="zh-cmn-Hans">
<head>
    <title>{{ .Test.Name }} {{ .I18n.Title.Report }}</title>
    <meta charset="UTF-8">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.8.1/font/bootstrap-icons.css">
    <script src="https://code.jquery.com/jquery-3.6.0.slim.min.js" integrity="sha256-u7e5khyithlIdTpu22PHhENmPcRdFiHRjhAuHcs05RI=" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/echarts@5.3.2/dist/echarts.min.js" integrity="sha256-7rldQObjnoCubPizkatB4UZ0sCQzu2ePgyGSUcVN70E=" crossorigin="anonymous"></script>

    {{ .Customize.Font.Style }}
    <style>
        body {
            font-family: {{ .Customize.Font.Body }};
        }
        pre, code {
            font-family: {{ .Customize.Font.Code }};
        }
    </style>

    <script>
    var yAxisLabelFormatter = {
        byte: (b) => {
          const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
          let l = 0, n = parseInt(b, 10) || 0;
          while(n >= 1024 && ++l){
              n = n/1024;
          }
          return(n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l]);
        },
        bit: (b) => {
          const units = ['b', 'Kb', 'Mb', 'Gb', 'Tb', 'Pb', 'Eb', 'Zb', 'Yb'];
          let l = 0, n = parseInt(b, 10) || 0;
          while(n >= 1024 && ++l){
              n = n/1024;
          }
          return(n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l]);
        },
        percent: (v) => {
            return v + "%";
        },
        times: (v) => {
          const units = ['', 'K', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y'];
          let l = 0, n = parseInt(v, 10) || 0;
          while(n >= 1024 && ++l){
              n = n/1024;
          }
          return(n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l]);
        }
    }

    </script>

    {{ .Customize.Extra.Head }}
</head>

<body>
    {{ .Customize.Extra.BodyHeader }}
    <div class="container">
        <div class="row justify-content-md-center">
            <div class="col-lg-10 col-md-12">
            {{ RenderTest .Test "test" }}
            </div>
        </div>
    </div>
    {{ .Customize.Extra.BodyFooter }}
</body>
<script>
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
      return new bootstrap.Tooltip(tooltipTriggerEl)
    })
</script>
</html>
`

var testTplStr = `
<div class="col-md-12" id={{ .Name }}>
    {{ if .Test.IsErr }}
    <div class="card my-{{ .Customize.Padding.Y }} border-danger">
        <h5 class="card-header text-white bg-danger">{{ .I18n.Title.Test }} {{ .Test.Name }} {{ .I18n.Status.Fail }}</h5>
    {{ else }}
    <div class="card my-{{ .Customize.Padding.Y }} border-success">
        <h5 class="card-header text-white bg-success">{{ .I18n.Title.Test }} {{ .Test.Name }} {{ .I18n.Status.Succ }}</h5>
    {{ end }}

        {{ if .Test.IsErr }}
        <div class="card-header text-white bg-danger"><span class="fw-bolder">{{ .I18n.Test.Err }}</span></div>
        <div class="card-body"><pre>{{ .Test.Err }}</pre></div>
        {{ end }}

        {{ if .Test.Description }}
        <div class="card-header justify-content-between d-flex"><span class="fw-bolder">{{ .I18n.TestHeader.Description }}</span></div>
        <div class="card-body">{{ Markdown .Test.Description }}</div>
        {{ end }}

        {{ if .Test.Plans }}
        <div class="card-header justify-content-between d-flex">
            <span class="fw-bolder">{{ .I18n.Title.Plan }}</span>
        </div>
        <ul class="list-group list-group-flush" id="{{ .Name }}-plan">
            {{ range $idx, $plan := .Test.Plans }}
            <li class="list-group-item px-{{ .Customize.Padding.X }} py-{{ .Customize.Padding.Y }} plan">
                {{ RenderPlan $plan (printf "%s-plan-%d" .Name $idx) }}
            </li>
            {{ end }}
        </ul>
        {{ end }}
    </div>
</div>
`

var planTplStr = `
<a class="card-title btn d-flex justify-content-between align-items-center" data-bs-toggle="collapse" href="#{{ .Name }}" role="button" aria-expanded="false" aria-controls="{{ .Name }}">
    {{ .Plan.Name }}
</a>
<div class="card collapse show" id="{{ .Name }}">
    {{ if plan.is_err }}
    <div class="card border-danger">
    {{ else }}
    <div class="card border-success">
    {{ end }}
    
        {{ if .Plan.Description }}
        <div class="card-header"><span class="fw-bolder">{{ I18n.Title.Description }}</span></div>
        <div class="card-body">{{ markdown(.Plan.Description) }}</div>
        {{ end }}
    
        {{ if .Plan.Command }}
        <div class="card-header"><span class="fw-bolder">{{ .I18n.Title.Command }}</span></div>
        <div class="card-body">
            <div class="float-end">
                <button type="button" class="btn btn-sm py-0" onclick="copyToClipboard('{{ .Name }}-command')"
                    data-bs-toggle="tooltip" data-bs-placement="top" title="{{ i18n.toolTips.copy }}">
                    <i class="bi-clipboard"></i>
                </button>
            </div>
            <span id="{{ .Name }}-command">{{ .Plan.Command }}</span>
        </div>
        {{ end }}
    
        {{ if plan.unit_groups }}
        <ul class="list-group list-group-flush">
            {{ for UnitGroup in .Plan.UnitGroups }}
            <li class="list-group-item px-{{ .Customize.Padding.X }} py-{{ .Customize.Padding.Y }}">
                {{ RenderUnitGroup .UnitGroup, '{}-group-{}'.format(name, loop.index0) }}
            </li>
            {{ endfor }}
        </ul>
        {{ end }}
    </div>
</div>
`

var unitGroupTplStr = `
<div class="card" id="{{ .Name }}">
    {{ if .UnitGroup.IsErr }}<div class="card border-danger">{{ else }}<div class="card border-success">{{ end }}

    <div class="card-header justify-content-between d-flex">
        <span class="fw-bolder">{{ .I18n.Title.Summary }} No.{{ group.idx + 1 }}</span>
        <span>
            {{ if .UnitGroup.Seconds }}
            <span class="badge bg-success rounded-pill">{{ .UnitGroup.Seconds }}s</span>
            {{ end }}
            {{ if .UnitGroup.Times }}
            <span class="badge bg-success rounded-pill">{{ .UnitGroup.Times }}</span>
            {{ end }}
        </span>
    </div>
    <div class="card-body">
        <table class="table table-striped">
            <thead>
                <tr class="text-center">
                    <th>{{ .I18n.Title.Unit }}</th>
                    <th>{{ .I18n.Title.Parallel }}</th>
                    <th>{{ .I18n.Title.Total }}</th>
                    <th>{{ .I18n.Title.Rate }}</th>
                    <th>{{ .I18n.Title.QPS }}</th>
                    <th>{{ .I18n.Title.ResTime }}</th>
                    {{ for q in .UnitGroup.Quantile }}
                    <th>{{ .I18n.Title.QuantileShort }}{{ q }}</th>
                    {{ endfor }}
                </tr>
            </thead>
            <tbody>
                {{ for unit in .UnitGroup.Units }}
                <tr class="text-center">
                    <td>{{ unit.name }}</td>
                    <td>{{ unit.parallel }}</td>
                    <td>{{ unit.total }}</td>
                    <td>{{ int(unit.rate * 10000) / 100 }}%</td>
                    <td>{{ int(unit.qps) }}</td>
                    <td>{{ format_timedelta(unit.res_time) }}</td>
                    {{ for q in .UnitGroup.Quantile }}
                    <td>{{ format_timedelta(unit.quantile[q]) }}</td>
                    {{ endfor }}
                </tr>
                {{ endfor }}
            </tbody>
        </table>
    </div>

    <div class="card-body d-flex justify-content-center">
        <div  class="col-md-12" id="{{ '{}-unit-code'.format(name) }}" style="height: 300px;"></div>
        <script>
            echarts.init(document.getElementById("{{ '{}-unit-code'.format(name) }}")).setOption({
              title: {
                text: "{{ .I18n.Title.Code }}",
                left: "center",
              },
              textStyle: {
                fontFamily: "{{ .Customize.Font.Echarts }}",
              },
              tooltip: {
                trigger: "item"
              },
              toolbox: {
                feature: {
                  saveAsImage: {
                    title: "{{ .I18n.Tooltip.Save }}"
                  }
                }
              },
              series: [
                {{ for unit in .UnitGroup.Units }}
                {
                  name: "{{ unit.name }}",
                  type: "pie",
                  radius: ['{{ (70 / loop.length) * loop.index0 + 15 }}%', '{{ (70 / loop.length) * loop.index + 10 }}%'],
                  avoidLabelOverlap: false,
                  label: {
                    show: false,
                    position: 'center'
                  },
                  emphasis: {
                    label: {
                      show: true,
                      fontSize: '20',
                      fontWeight: 'bold'
                    }
                  },
                  labelLine: {
                    show: false
                  },
                  data: {{ json.dumps(dict_to_items(unit.code)) }}
                },
                {{ endfor }}
              ]
            });
        </script>
    </div>

    <div class="card-body d-flex justify-content-center">
        <div class="col-md-12" id="{{ '{}-unit-qps'.format(name) }}" style="height: 300px;"></div>
        <script>
            echarts.init(document.getElementById("{{ '{}-unit-qps'.format(name) }}")).setOption({
              title: {
                text: "{{ .I18n.Title.QPS }}",
                left: "center",
              },
              textStyle: {
                fontFamily: "{{ .Customize.Font.Echarts }}",
              },
              tooltip: {
                trigger: 'axis',
                show: true,
                axisPointer: {
                    type: "cross"
                }
              },
              toolbox: {
                feature: {
                  saveAsImage: {
                    title: "{{ .I18n.Tooltip.Save }}"
                  }
                }
              },
              xAxis: {
                type: "time",
              },
              yAxis: {
                type: "value",
              },
              series: [
                {{ for unit in .UnitGroup.Units }}
                {
                  name: "{{ unit.name }}",
                  type: "line",
                  smooth: true,
                  symbol: "none",
                  areaStyle: {},
                  data: {{ json.dumps(unit_stage_serial(unit, "qps")) }}
                },
                {{ endfor }}
              ]
            });
        </script>
    </div>

    <div class="card-body d-flex justify-content-center">
        <div class="col-md-12" id="{{ '{}-unit-rate'.format(name) }}" style="height: 300px;"></div>
        <script>
            echarts.init(document.getElementById("{{ '{}-unit-rate'.format(name) }}")).setOption({
              title: {
                text: "{{ .I18n.Title.Rate }}",
                left: "center",
              },
              textStyle: {
                fontFamily: "{{ .Customize.Font.Echarts }}",
              },
              tooltip: {
                trigger: 'axis',
                show: true,
                axisPointer: {
                    type: "cross"
                }
              },
              toolbox: {
                feature: {
                  saveAsImage: {
                    title: "{{ .I18n.Tooltip.save }}"
                  }
                }
              },
              xAxis: {
                type: "time",
                boundaryGap: false
              },
              yAxis: {
                type: "value",
                axisLabel: {
                  formatter: yAxisLabelFormatter["percent"],
                }
              },
              series: [
                {{ for unit in .UnitGroup.Units }}
                {
                  name: "{{ unit.name }}",
                  type: "line",
                  smooth: true,
                  symbol: "none",
                  areaStyle: {},
                  data: {{ json.dumps(unit_stage_serial(unit, "rate", "percent")) }}
                },
                {{ endfor }}
              ]
            });
        </script>
    </div>
    
    {{ for mname, monitor in group.monitor.items() }}
    <div class="card-header justify-content-between d-flex"><span class="fw-bolder">{{ .I18n.Title.monitor }}-{{ mname }}</span></div>
    {{ for metric_name, stat in monitor["stat"].items() }}
    <div class="card-body d-flex justify-content-center">
        <div class="col-md-12" id="{{ '{}-monitor-{}-{}'.format(name, mname, metric_name) }}" style="height: 300px;"></div>
        <script>
            echarts.init(document.getElementById("{{ '{}-monitor-{}-{}'.format(name, mname, metric_name) }}")).setOption({
              title: {
                text: "{{ metric_name }}",
                left: "center",
              },
              textStyle: {
                fontFamily: "{{ .Customize.Font.Echarts }}",
              },
              tooltip: {
                trigger: 'axis',
                show: true,
                axisPointer: {
                    type: "cross"
                }
              },
              toolbox: {
                feature: {
                  saveAsImage: {
                    title: "{{ .I18n.Tooltip.save }}"
                  }
                }
              },
              xAxis: {
                type: "time",
                boundaryGap: false
              },
              yAxis: {
                type: "value",
                axisLabel: {
                  formatter: yAxisLabelFormatter["{{ monitor["unit"][metric_name] }}"],
                }
              },
              series: [
                {
                  name: "{{ metric_name }}",
                  type: "line",
                  smooth: true,
                  symbol: "none",
                  areaStyle: {},
                  data: {{ json.dumps(monitor_serial(stat, "value")) }}
                },
              ]
            });
        </script>
    </div>
    {{ endfor }}
    {{ endfor }}
</div>
`

type HtmlReporter struct {
	options *HtmlReporterOptions

	i18n *i18n.I18n

	reportTpl    *template.Template
	testTpl      *template.Template
	planTpl      *template.Template
	unitGroupTpl *template.Template
}

func (r *HtmlReporter) Report(test *stat.TestStat) string {
	var buf bytes.Buffer
	if err := r.reportTpl.Execute(&buf, map[string]interface{}{
		"Test":      test,
		"Customize": r.options,
		"I18n":      r.i18n,
	}); err != nil {
		return fmt.Sprintf("%+v", errors.Wrap(err, "r.reportTpl.Execute failed"))
	}

	return buf.String()
}

func (r *HtmlReporter) RenderTest(test *stat.TestStat, name string) string {
	var buf bytes.Buffer
	if err := r.reportTpl.Execute(&buf, map[string]interface{}{
		"Name":      name,
		"Test":      test,
		"Customize": r.options,
		"I18n":      r.i18n,
	}); err != nil {
		return fmt.Sprintf("%+v", errors.Wrap(err, "r.reportTpl.Execute failed"))
	}

	return buf.String()
}

func (r *HtmlReporter) RenderPlan(plan *stat.PlanStat, name string) string {
	var buf bytes.Buffer
	if err := r.planTpl.Execute(&buf, map[string]interface{}{
		"Name":      name,
		"Plan":      plan,
		"Customize": r.options,
		"I18n":      r.i18n,
	}); err != nil {
		return fmt.Sprintf("%+v", errors.Wrap(err, "r.reportTpl.Execute failed"))
	}

	return buf.String()
}

func (r *HtmlReporter) RenderUnitGroup(unitGroup *stat.UnitGroupStat, name string) string {
	var buf bytes.Buffer
	if err := r.planTpl.Execute(&buf, map[string]interface{}{
		"Name":      name,
		"UnitGroup": unitGroupTplStr,
		"Customize": r.options,
		"I18n":      r.i18n,
	}); err != nil {
		return fmt.Sprintf("%+v", errors.Wrap(err, "r.reportTpl.Execute failed"))
	}

	return buf.String()
}
