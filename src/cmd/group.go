package cmd

import (
	"bfv-bot/bot/group"
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/model/common/req"
	"bfv-bot/model/common/resp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func init() {
	groupCommandMap["ban"] = ban
	groupCommandMap["removeban"] = removeban
}

func cx(msg *req.MsgData, c *gin.Context, _ string, value string) {

	path, err := utils.QueryAndStore(value, 1)
	if err != nil {
		global.GLog.Error("utils.QueryAndStore, 1", zap.String("name", value), zap.Error(err))
		group.SendGroupMsg(msg.GroupID, "查询失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	global.GLog.Info("file:///" + path)

	// 直接发送图片消息到群聊，而不是返回特殊响应格式
	group.SendGroupImageMsg(msg.GroupID, "file:///"+path)
	resp.EmptyOk(c)
}

func ban(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	if !global.GConfig.QQBot.IsActiveAdminGroup(msg.GroupID) {
		resp.EmptyOk(c)
		return
	}
	group.SendGroupMsg(msg.GroupID, "[屏蔽] 已下线")
	resp.EmptyOk(c)
}

func platoon(msg *req.MsgData, c *gin.Context, _ string, value string) {
	err, result := utils.GetJoinPlatoonsByName(value)
	if err != nil {
		global.GLog.Error("utils.GetJoinPlatoonsByName",
			zap.String("name", value), zap.Error(err))
		group.SendGroupMsg(msg.GroupID, "查询失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	group.SendGroupMsg(msg.GroupID, result)
	resp.EmptyOk(c)
}

func banlog(msg *req.MsgData, c *gin.Context, _ string, value string) {
	err, result := utils.GetBanLog(value)
	if err != nil {
		global.GLog.Error("utils.GetBanLog",
			zap.String("name", value), zap.Error(err))
		group.SendGroupMsg(msg.GroupID, "查询失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	group.SendGroupMsg(msg.GroupID, result)
	resp.EmptyOk(c)
}

func removeban(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	if !global.GConfig.QQBot.IsActiveAdminGroup(msg.GroupID) {
		resp.EmptyOk(c)
		return
	}
	group.SendGroupMsg(msg.GroupID, "[解除屏蔽] 已下线")
	resp.EmptyOk(c)
}

func bind(msg *req.MsgData, c *gin.Context, _ string, value string) {
	err, data := utils.CheckPlayer(value)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, "绑定失败 "+err.Error())
		resp.EmptyOk(c)
		return
	}
	err = dbService.AddBind(msg.UserID, value, data.PID)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, "绑定失败 "+err.Error())
		resp.EmptyOk(c)
		return
	} else {
		group.SendGroupMsg(msg.GroupID, "绑定成功: "+data.PID)
		resp.EmptyOk(c)
		return
	}
}

func server(msg *req.MsgData, c *gin.Context, _ string, value string) {
	err, str := utils.GetBfvRobotServer(value, true)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, "查询失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	group.SendGroupMsg(msg.GroupID, str)
	resp.EmptyOk(c)
}

func data(msg *req.MsgData, c *gin.Context, _ string, value string) {

	path, err := utils.QueryAndStore(value, 2)
	if err != nil {
		global.GLog.Error("utils.QueryAndStore, 2", zap.String("name", value), zap.Error(err))
		group.SendGroupMsg(msg.GroupID, "查询失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	global.GLog.Info("file:///" + path)

	// 直接发送图片消息到群聊
	group.SendGroupImageMsg(msg.GroupID, "file:///"+path)
	resp.EmptyOk(c)
}

func task(msg *req.MsgData, c *gin.Context, _ string, value string) {

	offset, err := strconv.Atoi(value)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, "必须是数字")
		resp.EmptyOk(c)
		return
	}
	path, err := utils.GetTaskAndCache(offset)
	if err != nil {
		global.GLog.Error("utils.GetTaskAndCache", zap.String("value", value), zap.Error(err))
		group.SendGroupMsg(msg.GroupID, "获取失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	global.GLog.Info("file:///" + path)

	// 直接发送图片消息到群聊
	group.SendGroupImageMsg(msg.GroupID, "file:///"+path)
	resp.EmptyOk(c)
}

func playerlist(msg *req.MsgData, c *gin.Context, _ string, value string) {

	err, path := utils.GetPlayerList(value)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, "获取失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	global.GLog.Info("file:///" + path)

	// 直接发送图片消息到群聊
	group.SendGroupImageMsg(msg.GroupID, "file:///"+path)
	resp.EmptyOk(c)
}

func groupMember(msg *req.MsgData, c *gin.Context, _ string, value string) {

	err, s := utils.GerServerGroupMember(value)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, err.Error())
		resp.EmptyOk(c)
		return
	}
	group.SendGroupMsg(msg.GroupID, s)
	resp.EmptyOk(c)
}

func quickTask(msg *req.MsgData, c *gin.Context, key string) {
	task(msg, c, key, "0")
}

func ShortCommandFunction(msg *req.MsgData, c *gin.Context, command string) {
	err, name := dbService.GetBindName(msg.UserID)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, "快捷查询失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	groupCommandFunction, groupCommandOk := groupCommandMap[command]
	if groupCommandOk {
		groupCommandFunction(msg, c, command, name)
	}
}

func getGroupServerInfo(msg *req.MsgData, c *gin.Context, _ string) {
	err, result := utils.GetBfvRobotServer(global.GConfig.Bfv.GroupUniName, false)
	if err == nil {
		// 直接发送文本消息到群聊
		group.SendGroupMsg(msg.GroupID, result)
		resp.EmptyOk(c)
		return
	} else {
		// 直接发送错误消息到群聊
		group.SendGroupMsg(msg.GroupID, err.Error())
		resp.EmptyOk(c)
		return
	}
}

func quickCx(msg *req.MsgData, c *gin.Context, _ string, value string) {
	err, data := utils.CheckPlayer(value)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, "查询失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}

	err, finalMsg := utils.GetBaseInfoAndStatusByName(&data)
	if err != nil {
		group.SendGroupMsg(msg.GroupID, "查询失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	group.SendGroupMsg(msg.GroupID, "玩家 ["+data.Name+"] 基础数据如下\n\n"+finalMsg)
	resp.EmptyOk(c)
}

func getGroupHelpInfo(msg *req.MsgData, c *gin.Context, _ string) {
	var builder strings.Builder

	if len(global.GConfig.QQBot.CustomCommandKey.Banlog) != 0 {
		builder.WriteString("玩家屏蔽记录: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Banlog, "/") + "=<id>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Server) != 0 {
		builder.WriteString("服务器查询: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Server, "/") + "=<name>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Data) != 0 {
		builder.WriteString("完整数据查询: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Data, "/") + "=<name>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Task) != 0 {
		builder.WriteString("周任务查询: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Task, "/") + "=<上周: -1, 本周: 0, 下周: 1>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Playerlist) != 0 {
		builder.WriteString("服务器玩家列表查询: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Playerlist, "/") + "=<服务器>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.GroupMember) != 0 {
		builder.WriteString("查询服务器内群友: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.GroupMember, "/") + "=<服务器>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Cx) != 0 {
		builder.WriteString("战绩查询: " + strings.Join(global.GConfig.QQBot.CustomCommandKey.Cx, "/") + "=<id>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.C) != 0 {
		builder.WriteString("快捷查询: " + strings.Join(global.GConfig.QQBot.CustomCommandKey.C, "/") + "=<id>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Bind) != 0 {
		builder.WriteString("绑定玩家: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Bind, "/") + "=<id>\n")
	}

	if len(global.GConfig.QQBot.CustomCommandKey.Platoon) != 0 {
		builder.WriteString("加入的战排: " +
			strings.Join(global.GConfig.QQBot.CustomCommandKey.Platoon, "/") + "=<id>\n")
	}

	builder.WriteString("其他快捷指令: " + strings.Join(global.GConfig.QQBot.CustomCommandKey.GroupServer, "/") +
		"/" + global.GConfig.Bfv.GroupName)

	group.SendGroupMsg(msg.GroupID, builder.String())
	resp.EmptyOk(c)
}

func InitBanlogKey(key string) {
	groupCommandMap[key] = banlog
	groupShortCommandMap[key] = true
}

func InitCxKey(key string) {
	groupCommandMap[key] = cx
	groupShortCommandMap[key] = true
}

func InitCKey(key string) {
	groupShortCommandMap[key] = true
	groupCommandMap[key] = quickCx
}

func InitPlatoonKey(key string) {
	groupCommandMap[key] = platoon
	groupShortCommandMap[key] = true
}

func InitBindKey(key string) {
	groupCommandMap[key] = bind
}

func InitServerKey(key string) {
	groupCommandMap[key] = server
}

func InitDataKey(key string) {
	groupCommandMap[key] = data
	groupShortCommandMap[key] = true
}

func InitTaskKey(key string) {
	groupCommandMap[key] = task
}

func InitPlayerListKey(key string) {
	groupCommandMap[key] = playerlist
}

func InitGroupMemberKey(key string) {
	groupCommandMap[key] = groupMember
}

func InitHelpKey(key string) {
	groupQuickCommandMap[key] = getGroupHelpInfo
}

func InitGroupServerKey(key string) {
	groupQuickCommandMap[key] = getGroupServerInfo
}

func InitQuickTaskKey(key string) {
	groupQuickCommandMap[key] = quickTask
}

func GetGroupCommandFunc(key string) (func(*req.MsgData, *gin.Context, string, string), bool) {
	f, ok := groupCommandMap[key]
	return f, ok
}

func GetGroupShortCommandFunc(key string) (bool, bool) {
	f, ok := groupShortCommandMap[key]
	return f, ok
}

func GetGroupQuickCommandFunc(key string) (func(*req.MsgData, *gin.Context, string), bool) {
	f, ok := groupQuickCommandMap[key]
	return f, ok
}
