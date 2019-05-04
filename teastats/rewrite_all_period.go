package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 重写规则请求统计
type RewriteAllPeriodFilter struct {
	CounterFilter
}

func (this *RewriteAllPeriodFilter) Name() string {
	return "重写规则请求统计"
}

func (this *RewriteAllPeriodFilter) Codes() []string {
	return []string{
		"rewrite.all.second",
		"rewrite.all.minute",
		"rewrite.all.hour",
		"rewrite.all.day",
		"rewrite.all.week",
		"rewrite.all.month",
		"rewrite.all.year",
	}
}

func (this *RewriteAllPeriodFilter) Indexes() []string {
	return []string{"rewrite"}
}

func (this *RewriteAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *RewriteAllPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	if len(accessLog.RewriteId) == 0 {
		return
	}
	this.ApplyFilter(accessLog, map[string]string{
		"rewrite": accessLog.RewriteId,
	}, maps.Map{
		"count": 1,
	})
}

func (this *RewriteAllPeriodFilter) Stop() {
	this.StopFilter()
}
