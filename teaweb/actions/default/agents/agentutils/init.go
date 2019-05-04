package agentutils

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaconfigs/notices"
	"github.com/TeaWeb/code/teaweb/actions/default/notices/noticeutils"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/timers"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		// 检查Agent连通性
		checkConnecting()
	})
}

// 检查Agent连通性
func checkConnecting() {
	duration := 60 * time.Second
	maxDisconnections := 3
	timers.Loop(duration, func(looper *timers.Looper) {
		for _, agent := range agents.SharedAgents() {
			if !agent.On || !agent.CheckDisconnections {
				continue
			}

			runtimeAgent := FindAgentRuntime(agent)

			// 监控连通性
			_, isWaiting := CheckAgentIsWaiting(agent.Id)
			if !isWaiting {
				runtimeAgent.CountDisconnections++

				if runtimeAgent.CountDisconnections > 0 && runtimeAgent.CountDisconnections%maxDisconnections == 0 { // 失去连接 N 次则提醒
					sendDisconnectNotice(agent)
				}
			} else {
				hasDisconnections := runtimeAgent.CountDisconnections >= maxDisconnections
				runtimeAgent.CountDisconnections = 0

				// 发送恢复通知
				if hasDisconnections {
					sendConnectNotice(agent)
				}
			}
		}
	})
}

// 发送Agent失联通知
func sendDisconnectNotice(agent *agents.AgentConfig) {
	duration := 1 * time.Hour

	message := "Agent\"" + agent.Name + "（" + agent.Host + "）" + "\"失去连接"

	level := notices.NoticeLevelError
	t := time.Now()

	notice := notices.NewNotice()
	notice.Id = primitive.NewObjectID()
	notice.Agent.AgentId = agent.Id
	notice.Agent.Level = level
	notice.Message = message
	notice.SetTime(t)
	notice.Hash()

	// 同样的消息短时间内只发送一条
	if noticeutils.ExistNoticesWithHash(notice.MessageHash, map[string]interface{}{
		"agent.agentId": agent.Id,
		"agent.appId":   "",
		"agent.itemId":  "",
	}, duration) {
		return
	}

	err := noticeutils.NewNoticeQuery().Insert(notice)
	if err != nil {
		logs.Error(err)
	} else {
		// 通过媒介发送通知
		setting := notices.SharedNoticeSetting()
		fullMessage := "消息：" + message + "\n时间：" + timeutil.Format("Y-m-d H:i:s", t)
		linkNames := []string{}
		for _, l := range FindNoticeLinks(notice) {
			linkNames = append(linkNames, types.String(l["name"]))
		}
		if len(linkNames) > 0 {
			fullMessage += "\n位置：" + strings.Join(linkNames, "/")
		}

		// 查找分组，如果分组中有通知设置，则使用分组中的通知设置
		isNotified := false
		receiverIds := []string{}

		// 查找agent设置
		{
			receivers, found := agent.NoticeSetting[level]
			if found && len(receivers) > 0 {
				isNotified = true
				receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]失去连接", fullMessage, func(receiverId string, minutes int) int {
					return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
						"agent.agentId": agent.Id,
						"agent.appId":   "",
					}, minutes)
				})
			}
		}

		// 查找group设置
		if !isNotified {
			groupId := ""
			if len(agent.GroupIds) > 0 {
				groupId = agent.GroupIds[0]
			}
			group := agents.SharedGroupConfig().FindGroup(groupId)
			if group != nil {
				receivers, found := group.NoticeSetting[level]
				if found && len(receivers) > 0 {
					isNotified = true
					receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]失去连接", fullMessage, func(receiverId string, minutes int) int {
						return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
							"agent.agentId": agent.Id,
							"agent.appId":   "",
						}, minutes)
					})
				}
			}
		}

		// 默认通知媒介
		if !isNotified {
			receiverIds = setting.Notify(level, "["+agent.GroupName()+"]["+agent.Name+"]失去连接", fullMessage, func(receiverId string, minutes int) int {
				return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
					"agent.agentId": agent.Id,
					"agent.appId":   "",
				}, minutes)
			})
		}

		if len(receiverIds) > 0 {
			noticeutils.UpdateNoticeReceivers(notice.Id, receiverIds)
		}
	}
}

// 发送Agent连接通知
func sendConnectNotice(agent *agents.AgentConfig) {
	message := "Agent\"" + agent.Name + "（" + agent.Host + "）" + "\"已恢复连接"

	level := notices.NoticeLevelSuccess
	t := time.Now()

	notice := notices.NewNotice()
	notice.Id = primitive.NewObjectID()
	notice.Agent.AgentId = agent.Id
	notice.Agent.Level = level
	notice.Message = message
	notice.SetTime(t)
	notice.Hash()
	err := noticeutils.NewNoticeQuery().Insert(notice)
	if err != nil {
		logs.Error(err)
	} else {
		// 通过媒介发送通知
		setting := notices.SharedNoticeSetting()
		fullMessage := "消息：" + message + "\n时间：" + timeutil.Format("Y-m-d H:i:s", t)
		linkNames := []string{}
		for _, l := range FindNoticeLinks(notice) {
			linkNames = append(linkNames, types.String(l["name"]))
		}
		if len(linkNames) > 0 {
			fullMessage += "\n位置：" + strings.Join(linkNames, "/")
		}

		// 查找分组，如果分组中有通知设置，则使用分组中的通知设置
		isNotified := false
		receiverIds := []string{}

		// 查找agent设置
		{
			receivers, found := agent.NoticeSetting[level]
			if found && len(receivers) > 0 {
				isNotified = true
				receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]恢复连接", fullMessage, func(receiverId string, minutes int) int {
					return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
						"agent.agentId": agent.Id,
						"agent.appId":   "",
					}, minutes)
				})
			}
		}

		// 查找group设置
		if !isNotified {
			groupId := ""
			if len(agent.GroupIds) > 0 {
				groupId = agent.GroupIds[0]
			}
			group := agents.SharedGroupConfig().FindGroup(groupId)
			if group != nil {
				receivers, found := group.NoticeSetting[level]
				if found && len(receivers) > 0 {
					isNotified = true
					receiverIds = setting.NotifyReceivers(level, receivers, "["+agent.GroupName()+"]["+agent.Name+"]恢复连接", fullMessage, func(receiverId string, minutes int) int {
						return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
							"agent.agentId": agent.Id,
							"agent.appId":   "",
						}, minutes)
					})
				}
			}
		}

		// 默认通知媒介
		if !isNotified {
			receiverIds = setting.Notify(level, "["+agent.GroupName()+"]["+agent.Name+"]恢复连接", fullMessage, func(receiverId string, minutes int) int {
				return noticeutils.CountReceivedNotices(receiverId, map[string]interface{}{
					"agent.agentId": agent.Id,
					"agent.appId":   "",
				}, minutes)
			})
		}

		if len(receiverIds) > 0 {
			noticeutils.UpdateNoticeReceivers(notice.Id, receiverIds)
		}
	}
}
