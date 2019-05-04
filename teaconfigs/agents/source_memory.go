package agents

import (
	"github.com/TeaWeb/code/teaconfigs/forms"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaconfigs/widgets"
	"github.com/iwind/TeaGo/maps"
)

// 内存使用量
type MemorySource struct {
	Source `yaml:",inline"`
}

// 获取新对象
func NewMemorySource() *MemorySource {
	return &MemorySource{}
}

// 名称
func (this *MemorySource) Name() string {
	return "内存"
}

// 代号
func (this *MemorySource) Code() string {
	return "memory"
}

// 描述
func (this *MemorySource) Description() string {
	return "内存使用量等信息"
}

// 表单信息
func (this *MemorySource) Form() *forms.Form {
	form := forms.NewForm(this.Code())
	return form
}

// 变量
func (this *MemorySource) Variables() []*SourceVariable {
	return []*SourceVariable{
		{
			Code:        "usage.virtualPercent",
			Description: "Virtual内存使用量百分比",
		},
		{
			Code:        "usage.virtualUsed",
			Description: "Virtual内存使用量（G）",
		},
		{
			Code:        "usage.virtualTotal",
			Description: "Virtual总内存容量（G）",
		},
		{
			Code:        "usage.virtualFree",
			Description: "Free内存容量（G）",
		},
		{
			Code:        "usage.virtualWired",
			Description: "Wired内存",
		},
		{
			Code:        "usage.virtualBuffers",
			Description: "Buffers内存",
		},
		{
			Code:        "usage.virtualCached",
			Description: "Cached内存",
		},
		{
			Code:        "usage.swapPercent",
			Description: "Swap内存使用量百分比",
		},
		{
			Code:        "usage.swapUsed",
			Description: "Swap内存使用量（G）",
		},
		{
			Code:        "usage.swapTotal",
			Description: "Swap总内存容量（G）",
		},
		{
			Code:        "usage.swapFree",
			Description: "Swap Free内存容量（G）",
		},
	}
}

// 阈值
func (this *MemorySource) Thresholds() []*Threshold {
	result := []*Threshold{}

	{
		t := NewThreshold()
		t.Param = "${usage.virtualPercent}"
		t.Operator = ThresholdOperatorGte
		t.Value = "80"
		t.NoticeLevel = notices.NoticeLevelWarning
		result = append(result, t)
	}

	return result
}

// 图表
func (this *MemorySource) Charts() []*widgets.Chart {
	charts := []*widgets.Chart{}

	// chart
	{
		chart := widgets.NewChart()
		chart.Name = "内存使用量（%）"
		chart.Columns = 2
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `
var chart = new charts.LineChart();

var query = new values.Query();
query.limit(30)
var ones = query.desc().cache(60).findAll();
ones.reverse();

var lines = [];

{
	var line = new charts.Line();
	line.color = colors.ARRAY[0];
	line.isFilled = true;
	line.values = [];
	lines.push(line);
}

ones.$each(function (k, v) {
	lines[0].values.push(v.value.usage.virtualPercent);

	var minute = v.timeFormat.minute.substring(8);
	chart.labels.push(minute.substr(0, 2) + ":" + minute.substr(2, 2));
});

chart.addLines(lines);
chart.max = 100;
chart.render();
`,
		}
		charts = append(charts, chart)
	}

	{
		chart := widgets.NewChart()
		chart.Name = "当前内存使用量"
		chart.Columns = 1
		chart.Type = "javascript"
		chart.Options = maps.Map{
			"code": `
var chart = new charts.StackBarChart();

var latest = new values.Query().latest(1);
var hasWarning = false;
if (latest.length > 0) {
	hasWarning = (latest[0].value.usage.swapPercent > 50) || (latest[0].value.usage.virtualPercent > 80);
	chart.values = [ 
		[latest[0].value.usage.swapUsed, latest[0].value.usage.swapTotal - latest[0].value.usage.swapUsed],
		[latest[0].value.usage.virtualUsed, latest[0].value.usage.virtualTotal - latest[0].value.usage.virtualUsed]
	];
	chart.labels = [ "虚拟内存（" +  (Math.round(latest[0].value.usage.swapUsed * 10) / 10) + "G/" + Math.round(latest[0].value.usage.swapTotal) + "G"  + "）", "物理内存（" + (Math.round(latest[0].value.usage.virtualUsed * 10) / 10)+ "G/" + Math.round(latest[0].value.usage.virtualTotal)  + "G"  + "）"];
} else {
	chart.values = [ [0, 0], [0, 0] ];
	chart.labels = [ "虚拟内存", "物理内存" ];
}
if (hasWarning) {
	chart.colors = [ colors.RED, colors.GREEN ];
} else {
	chart.colors = [ colors.BROWN, colors.GREEN ];
}
chart.render();
`,
		}
		charts = append(charts, chart)
	}

	return charts
}
