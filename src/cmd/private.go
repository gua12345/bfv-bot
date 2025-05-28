package cmd

import (
	"bfv-bot/bot/group"
	"bfv-bot/bot/private"
	"bfv-bot/common/global"
	"bfv-bot/flow"
	"bfv-bot/model/common/req"
	"bfv-bot/model/common/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func init() {
	privateCommandMap["addblack"] = addblack
	privateCommandMap["removeblack"] = removeblack
	privateCommandMap["removecardcheck"] = removecardcheck
	privateCommandMap["addsensitive"] = addsensitive
	privateCommandMap["removesensitive"] = removesensitive
	privateCommandMap["addjoinblacklist"] = addjoinblacklist
	privateCommandMap["removejoinblacklist"] = removejoinblacklist
	privateCommandMap["bindtoken"] = bindtoken
	privateCommandMap["bindgameid"] = bindgameid
	privateCommandMap["op"] = op

	privateOpCommandMap["start"] = opStart
	privateOpCommandMap["stop"] = opStop

	privateOpCommandMap["start-broadcast"] = opStartBroadcast
	privateOpCommandMap["stop-broadcast"] = opStopBroadcast

	privateOpCommandMap["checknow"] = opChecknow
	privateOpCommandMap["gameid"] = opGameid
	privateOpCommandMap["token"] = opToken
	privateOpCommandMap["joinblacklist"] = opJoinBlackList
	privateOpCommandMap["deletejoinblacklist"] = opDeletejoinblacklist
	privateOpCommandMap["blacklist"] = opBlacklist
	privateOpCommandMap["sensitive"] = opSensitive
	privateOpCommandMap["grouplist"] = opGroupList
	privateOpCommandMap["getmsg"] = opGetMsg
	privateOpCommandMap["grouphistory"] = opGroupHistory

	privateQuickCommandMap["help"] = getPrivateHelpInfo
	privateQuickCommandMap[".help"] = getPrivateHelpInfo

}

func opStart(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	global.GConfig.Bfv.Active = true
	private.SendPrivateMsg(msg.UserID, "开始检测")
	resp.EmptyOk(c)
}

func opStop(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	global.GConfig.Bfv.Active = false
	global.GConfig.Bfv.ClearGameId()
	private.SendPrivateMsg(msg.UserID, "结束检测")
	resp.EmptyOk(c)
}

func opStartBroadcast(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	global.GConfig.QQBot.BotToBot.Enable = true
	private.SendPrivateMsg(msg.UserID, "开始喊话")
	resp.EmptyOk(c)
}

func opStopBroadcast(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	global.GConfig.QQBot.BotToBot.Enable = false
	private.SendPrivateMsg(msg.UserID, "结束喊话")
	resp.EmptyOk(c)
}

func opChecknow(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	cronService.CheckBlackListAndNotify()
	private.SendPrivateMsg(msg.UserID, "立即检测")
	resp.EmptyOk(c)
}

func opGameid(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	var builder strings.Builder
	for _, info := range global.GConfig.Bfv.Server {
		builder.WriteString(info.ServerName)
		builder.WriteString("\n")
		if info.GetGameId() == "" {
			builder.WriteString("无")
		} else {
			builder.WriteString(info.GetGameId())
		}
		builder.WriteString("\n")
	}
	resp.ReplyOk(c, builder.String())
}

func opToken(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	var builder strings.Builder
	for _, info := range global.GConfig.Bfv.Server {
		builder.WriteString(info.ServerName)
		builder.WriteString("\n")
		if info.GetToken() == "" {
			builder.WriteString("无")
		} else {
			builder.WriteString(info.GetToken())
		}
		builder.WriteString("\n")
	}
	resp.ReplyOk(c, builder.String())
}

func opJoinBlackList(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	list := dbService.QueryAllJoinBlackList()
	var builder strings.Builder
	builder.WriteString("加群黑名单\n")
	for key, value := range list {
		builder.WriteString(strconv.FormatInt(key, 10))
		builder.WriteString("\t")
		builder.WriteString(value)
		builder.WriteString("\n")
		builder.WriteString("\n")
	}
	finalStr := builder.String()
	if len(finalStr) > 0 {
		finalStr = finalStr[:len(finalStr)-1]
	}
	resp.ReplyOk(c, finalStr)
}

func opDeletejoinblacklist(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	err := dbService.DeleteAllJoinBlackList()
	if err != nil {
		resp.ReplyOk(c, "清空加群黑名单失败")
	} else {
		resp.ReplyOk(c, "清空加群黑名单成功")
	}
}

func opBlacklist(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	list := dbService.QueryAllBlackList()
	var builder strings.Builder
	builder.WriteString("黑名单\n")
	for key, value := range list {
		builder.WriteString("pid: ")
		builder.WriteString(key)
		builder.WriteString("\t")
		builder.WriteString("id: ")
		builder.WriteString(value.Name)
		builder.WriteString("\t")
		builder.WriteString("原因: ")
		builder.WriteString(value.Reason)
		builder.WriteString("\n")
	}
	resp.ReplyOk(c, builder.String())
}

func opSensitive(_ *req.MsgData, c *gin.Context, _ string, _ string) {
	list := dbService.SelectAllSensitive()
	var builder strings.Builder
	builder.WriteString("敏感词\n")
	for index, item := range list {
		builder.WriteString(strconv.Itoa(index + 1))
		builder.WriteString(". ")
		builder.WriteString(item)
		builder.WriteString("\n")
	}
	resp.ReplyOk(c, builder.String())
}

func addblack(msg *req.MsgData, c *gin.Context, _ string, value string) {

	flow.InitPrivateFlow(msg.UserID, msg.MessageID, flow.AddBlack, value)
	private.SendPrivateMsg(msg.UserID, "[添加黑名单] 请输入原因")

	resp.EmptyOk(c)
}

func removeblack(_ *req.MsgData, c *gin.Context, _ string, value string) {

	err := dbService.RemoveBlack(value)
	if err != nil {
		resp.ReplyOk(c, "移除失败")
	} else {
		resp.ReplyOk(c, fmt.Sprintf("黑名单用户 [%s] 移除成功", value))
	}
}

func removecardcheck(_ *req.MsgData, c *gin.Context, _ string, value string) {

	qq, _ := strconv.ParseInt(value, 10, 64)
	err := dbService.DeleteCardCheck(qq)
	if err != nil {
		resp.ReplyOk(c, "移除失败")
	} else {
		resp.ReplyOk(c, fmt.Sprintf("ID检测 [%s] 移除成功", value))
	}
}

func addsensitive(_ *req.MsgData, c *gin.Context, _ string, value string) {

	err := dbService.AddSensitive(value)
	if err != nil {
		resp.ReplyOk(c, "添加失败")
	} else {
		resp.ReplyOk(c, fmt.Sprintf("添加成功"))
		global.GSensitive.AddWord(value)
	}
}

func removesensitive(_ *req.MsgData, c *gin.Context, _ string, value string) {
	err := dbService.RemoveSensitive(value)
	if err != nil {
		resp.ReplyOk(c, "移除失败")
	} else {
		resp.ReplyOk(c, "移除成功, 重启生效")
	}

}

func addjoinblacklist(msg *req.MsgData, c *gin.Context, _ string, value string) {

	flow.InitPrivateFlow(msg.UserID, msg.MessageID, flow.AddJoinBlack, value)
	private.SendPrivateMsg(msg.UserID, "[添加加群黑名单] 请输入原因")

	resp.EmptyOk(c)
}

func removejoinblacklist(_ *req.MsgData, c *gin.Context, _ string, value string) {
	err := dbService.RemoveJoinBlackList(value)
	if err != nil {
		resp.ReplyOk(c, "移除失败")
	} else {
		resp.ReplyOk(c, fmt.Sprintf("[移除加群黑名单] [%s] 移除成功", value))
	}
}

func bindtoken(msg *req.MsgData, c *gin.Context, _ string, value string) {
	flow.InitPrivateFlow(msg.UserID, msg.MessageID, flow.BindToken, value)
	private.SendPrivateMsg(msg.UserID, "[绑定Token] 请输入服务器ID")

	resp.EmptyOk(c)
}

func bindgameid(msg *req.MsgData, c *gin.Context, _ string, value string) {
	flow.InitPrivateFlow(msg.UserID, msg.MessageID, flow.BindGameID, value)
	private.SendPrivateMsg(msg.UserID, "[绑定GameID] 请输入服务器ID")

	resp.EmptyOk(c)
}

func op(msg *req.MsgData, c *gin.Context, key string, value string) {
	function, ok := privateOpCommandMap[value]
	if ok {
		function(msg, c, key, value)
		return
	}
	resp.EmptyOk(c)
}

// 新增：获取群列表
func opGroupList(msg *req.MsgData, c *gin.Context, _ string, _ string) {
	err, groupList := group.GetGroupList(false)
	if err != nil {
		private.SendPrivateMsg(msg.UserID, "获取群列表失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	var builder strings.Builder
	builder.WriteString("群列表:\n")
	for _, groupInfo := range groupList {
		builder.WriteString(fmt.Sprintf("群号: %d, 群名: %s\n", groupInfo.GroupId, groupInfo.GroupName))
	}
	private.SendPrivateMsg(msg.UserID, builder.String())
	resp.EmptyOk(c)
}

// 新增：获取消息
func opGetMsg(msg *req.MsgData, c *gin.Context, _ string, value string) {
	messageId, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		private.SendPrivateMsg(msg.UserID, "消息ID必须是数字")
		resp.EmptyOk(c)
		return
	}
	err, msgData := group.GetMsg(messageId)
	if err != nil {
		private.SendPrivateMsg(msg.UserID, "获取消息失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("消息ID: %d\n", msgData.MessageId))
	builder.WriteString(fmt.Sprintf("消息类型: %s\n", msgData.MessageType))
	builder.WriteString(fmt.Sprintf("发送者: %s(%d)\n", msgData.Sender.Nickname, msgData.Sender.UserId))
	builder.WriteString(fmt.Sprintf("时间: %d\n", msgData.Time))
	private.SendPrivateMsg(msg.UserID, builder.String())
	resp.EmptyOk(c)
}

// 新增：获取群历史聊天记录
func opGroupHistory(msg *req.MsgData, c *gin.Context, _ string, value string) {
	// 解析参数: groupId,messageId,count
	parts := strings.Split(value, ",")
	if len(parts) < 2 {
		private.SendPrivateMsg(msg.UserID, "参数格式: 群号,消息ID[,数量]")
		resp.EmptyOk(c)
		return
	}

	groupId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		private.SendPrivateMsg(msg.UserID, "群号必须是数字")
		resp.EmptyOk(c)
		return
	}

	messageId := parts[1]
	count := 20 // 默认20条
	if len(parts) >= 3 {
		if c, err := strconv.Atoi(parts[2]); err == nil {
			count = c
		}
	}

	err, historyData := group.GetGroupMsgHistory(groupId, messageId, count)
	if err != nil {
		private.SendPrivateMsg(msg.UserID, "获取群历史记录失败: "+err.Error())
		resp.EmptyOk(c)
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("群%d的历史消息(%d条):\n", groupId, len(historyData.Messages)))
	for i, msg := range historyData.Messages {
		if i >= 10 { // 最多显示10条
			builder.WriteString("...(更多消息)\n")
			break
		}
		builder.WriteString(fmt.Sprintf("%s: %s\n", msg.Sender.Nickname, msg.RawMessage))
	}
	private.SendPrivateMsg(msg.UserID, builder.String())
	resp.EmptyOk(c)
}

func getPrivateHelpInfo(msg *req.MsgData, c *gin.Context, _ string) {
	var builder strings.Builder
	builder.WriteString("绑定token: bindtoken=<token>\n")
	builder.WriteString("绑定gameid: bindgameid=<gameid>\n")
	builder.WriteString("添加黑名单: addblack=<id>\n")
	builder.WriteString("移除黑名单: removeblack=<id>\n")
	builder.WriteString("移除id检测: removecardcheck=<qq>\n")
	builder.WriteString("添加敏感词: addsensitive=<id>\n")
	builder.WriteString("移除敏感词: removesensitive=<id>\n")
	builder.WriteString("添加加群黑名单: addjoinblacklist=<qq>\n")
	builder.WriteString("移除加群黑名单: removejoinblacklist=<qq>\n")
	builder.WriteString("获取游戏id: op=gameid\n")
	builder.WriteString("获取服务器token: op=token\n")
	builder.WriteString("开始检测黑名单: op=start\n")
	builder.WriteString("停止检测黑名单: op=stop\n")
	builder.WriteString("开始喊话: op=start-broadcast\n")
	builder.WriteString("停止喊话: op=stop-broadcast\n")
	builder.WriteString("立即检测黑名单: op=checknow\n")
	builder.WriteString("清空加群黑名单: op=deletejoinblacklist\n")
	builder.WriteString("加群黑名单列表: op=joinblacklist\n")
	builder.WriteString("敏感词列表: op=sensitive\n")
	builder.WriteString("黑名单列表: op=blacklist\n")
	builder.WriteString("获取群列表: op=grouplist\n")
	builder.WriteString("获取消息: op=getmsg=<消息ID>\n")
	builder.WriteString("获取群历史: op=grouphistory=<群号,消息ID[,数量]>")
	private.SendPrivateMsg(msg.UserID, builder.String())
	resp.EmptyOk(c)
}

func GetPrivateCommandFunc(key string) (func(*req.MsgData, *gin.Context, string, string), bool) {
	f, ok := privateCommandMap[key]
	return f, ok
}

func GetPrivateQuickCommandFunc(key string) (func(*req.MsgData, *gin.Context, string), bool) {
	f, ok := privateQuickCommandMap[key]
	return f, ok
}
