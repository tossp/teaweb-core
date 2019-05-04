package teastats

import (
	"github.com/TeaWeb/code/tealogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net/url"
	"strings"
)

// 来源域名统计
type RefererDomainPeriodFilter struct {
	CounterFilter
}

func (this *RefererDomainPeriodFilter) Name() string {
	return "来源域名统计"
}

func (this *RefererDomainPeriodFilter) Codes() []string {
	return []string{
		"referer.domain.second",
		"referer.domain.minute",
		"referer.domain.hour",
		"referer.domain.day",
		"referer.domain.week",
		"referer.domain.month",
		"referer.domain.year",
	}
}

func (this *RefererDomainPeriodFilter) Indexes() []string {
	return []string{"domain"}
}

func (this *RefererDomainPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *RefererDomainPeriodFilter) Filter(accessLog *tealogs.AccessLog) {
	referer := accessLog.GetHeader("Referer")
	if len(referer) == 0 {
		return
	}
	uri, err := url.Parse(referer)
	if err != nil {
		return
	}
	domain := uri.Host

	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}

	this.ApplyFilter(accessLog, map[string]string{
		"domain": domain,
	}, maps.Map{
		"count": 1,
	})
}

func (this *RefererDomainPeriodFilter) Stop() {
	this.StopFilter()
}
