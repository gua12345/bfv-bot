package api

import (
	"bfv-bot/bot/group"
	"bfv-bot/bot/private"
	"bfv-bot/cmd"
	"bfv-bot/common/cons"
	"bfv-bot/common/des"
	"bfv-bot/common/global"
	"bfv-bot/common/utils"
	"bfv-bot/flow"
	"bfv-bot/model/common/req"
	"bfv-bot/model/common/resp"
	"context"
	"fmt"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type EventApi struct{}

// Post 事件
func (a *EventApi) Post(c *gin.Context) {
	rawData, err := c.GetRawData()
	if err != nil {
		resp.EmptyOk(c)
		return
	}
	var data req.BaseData
	err = des.ByteToStruct(rawData, &data)
	if err != nil {
		global.GLog.Error("utils.ByteToStruct", zap.Error(err))
		resp.EmptyOk(c)
		return
	}

	// 通知类型是消息
	if data.PostType == "message" {
		var msg req.MsgData
		err = des.ByteToStruct(rawData, &msg)
		if err != nil {
			global.GLog.Error("utils.ByteToStruct", zap.Error(err))
			resp.EmptyOk(c)
			return
		}

		// 消息格式检查：新接口没有message_format字段，旧接口需要是array格式
		// 如果有message_format字段且不是array，则不处理（兼容旧接口）
		if msg.MessageFormat != "" && msg.MessageFormat != "array" {
			global.GLog.Info("消息格式不支持", zap.String("message_format", msg.MessageFormat))
			resp.EmptyOk(c)
			return
		}
		// 群聊指令
		if msg.MessageType == "group" {

			// 没有启用的群的消息不处理
			if !global.GConfig.QQBot.IsActiveGroup(msg.GroupID) {
				global.GLog.Info("群组未激活",
					zap.Int64("群ID", msg.GroupID),
					zap.Any("激活群列表", global.GConfig.QQBot.ActiveGroup))
				resp.EmptyOk(c)
				return
			}

			global.GLog.Info("群组已激活，开始处理消息", zap.Int64("群ID", msg.GroupID))
			// 非文本消息不处理
			if msg.Message[0].Type == "text" {
				// 是否在流程中
				haveNext := flow.DoGroupNextStep(&msg)
				if haveNext {
					resp.EmptyOk(c)
					return
				}
				command := msg.Message[0].Data.Text

				// 命令匹配顺序
				// 1. 完整命令 (cx=id)
				// 2. 短命令 (cx/banlog)
				// 3. 快捷命令 (查服/任务)
				// other: 敏感词检测
				key, value := utils.GetCommandKeyValue(command)

				// 添加日志输出
                global.GLog.Info("已接收到命令",
				zap.String("command", command),
				zap.String("key", key),
				zap.String("value", value))

				groupCommandFunction, groupCommandOk := cmd.GetGroupCommandFunc(key)
				_, shortCommandOk := cmd.GetGroupShortCommandFunc(command)

				groupQuickCommandFunction, groupQuickCommandOk := cmd.GetGroupQuickCommandFunc(command)

				// 添加更多调试信息
				global.GLog.Info("命令匹配结果",
					zap.Bool("groupCommandOk", groupCommandOk),
					zap.Bool("shortCommandOk", shortCommandOk),
					zap.Bool("groupQuickCommandOk", groupQuickCommandOk))

				if groupCommandOk {
					global.GLog.Info("执行群组命令", zap.String("key", key), zap.String("value", value))
					groupCommandFunction(&msg, c, key, value)
					return
				} else if shortCommandOk {
					global.GLog.Info("执行短命令", zap.String("command", command))
					cmd.ShortCommandFunction(&msg, c, command)
					return
				} else if groupQuickCommandOk {
					global.GLog.Info("执行快捷命令", zap.String("command", command))
					groupQuickCommandFunction(&msg, c, command)
					return
				} else {
					global.GLog.Info("没有匹配到任何命令", zap.String("key", key), zap.String("command", command))
				}

				// 敏感词检测
				match, _ := global.GSensitive.Find(msg.RawMessage)
				if match {
					group.DeleteMsg(msg.MessageID)
					group.SendGroupMsg(msg.GroupID, "敏感话题, 不要讨论了")
				}

			} else if msg.Message[0].Type == "at" && msg.Message[0].Data.Qq == global.GConfig.QQBot.Qq {
				if global.GConfig.Ai.Enable {
					// ai回复
					seed := time.Now().UnixNano()
					rand.New(rand.NewSource(seed))

					num := rand.Intn(100)
					if num >= 10 {
						resp.EmptyOk(c)
						return
					}
					if len(msg.Message) == 1 {
						resp.ReplyOk(c, "？")
						return
					} else if len(msg.Message) == 2 {
						if msg.Message[1].Type == "text" {
							resp.ReplyOk(c, aiMsg(msg.Message[1].Data.Text))
							return
						}
					}
				}
			}

		} else if msg.MessageType == "private" {

			// 私聊管理指令 必须是管理员才能操作
			if !global.GConfig.QQBot.IsActiveAdminQq(msg.Sender.UserID) {
				resp.EmptyOk(c)
				return
			}

			if msg.Message[0].Type == "text" {
				// 判断是否存在已有流程
				haveNext := flow.DoPrivateNextStep(&msg)
				if haveNext {
					resp.EmptyOk(c)
					return
				}
				command := msg.Message[0].Data.Text

				key, value := utils.GetCommandKeyValue(command)
				function, ok := cmd.GetPrivateCommandFunc(key)

				privateQuickCommandFunction, privateQuickCommandOk := cmd.GetPrivateQuickCommandFunc(command)

				if ok {
					function(&msg, c, key, value)
					return
				} else if privateQuickCommandOk {
					privateQuickCommandFunction(&msg, c, command)
					return
				}
			}

		}
	} else if data.PostType == "request" {
		var msg req.AddGroupData
		err := des.ByteToStruct(rawData, &msg)
		if err != nil {
			global.GLog.Error("utils.ByteToStruct", zap.Error(err))
			resp.EmptyOk(c)
			return
		}
		// 当前群得是激活状态
		// 加群处理逻辑
		if global.GConfig.QQBot.IsActiveGroup(msg.GroupID) && msg.RequestType == "group" && msg.SubType == "add" {
			var approve bool = true
			var reason string = ""

			// 优先判断黑名单
			value, ok := global.GJoinBlackListMap[msg.UserID]
			if ok {
				approve = false
				reason = fmt.Sprintf("黑名单, 原因: %s", value)
			} else if global.GConfig.QQBot.EnableRejectFullGroupJoinRequest {
				// 检查群是否已满
				err, groupInfo := group.GetGroupInfo(msg.GroupID, true)
				if err != nil {
					global.GLog.Error("GetGroupInfo", zap.Error(err))
					// 如果获取群信息失败，继续正常流程
				} else if groupInfo.MemberCount >= groupInfo.MaxMemberCount {
					approve = false
					reason = "群聊已满，暂时无法加入"
					global.GLog.Info("拒绝加群申请：群聊已满",
						zap.Int64("群ID", msg.GroupID),
						zap.Int64("申请人", msg.UserID),
						zap.Int("当前人数", groupInfo.MemberCount),
						zap.Int("最大人数", groupInfo.MaxMemberCount))
				}
			}

			if approve {
				match := GroupAnswerReg.FindStringSubmatch(msg.Comment)
				if len(match) > 1 {
					name := strings.TrimSpace(match[1])
					if name == "" {
						approve = false
						reason = "未提供id"
					} else {
						err, data := utils.CheckPlayer(url.QueryEscape(name))

						if err != nil || data.PID == "" {
							if global.GConfig.QQBot.EnableRejectJoinRequest {
								approve = false
								if err != nil {
									if err.Error() == cons.PlayerNotFound {
										reason = "未能确认你提供的id"
									} else {
										reason = "其他异常: " + err.Error()
									}
								} else {
									reason = "pid获取失败"
								}
							} else {
								approve = true

								_ = global.GPool.Submit(func() {
									time.Sleep(1 * time.Second)
									// 欢迎信息
									group.SendAtGroupMsg(msg.GroupID, msg.UserID, global.GConfig.QQBot.WelcomeMsg)
									global.GLog.Error("utils.CheckPlayer", zap.Error(err))
									if err.Error() == cons.PlayerNotFound {

										content := " 机器人无法确认你提供的ID: [" + name + "]，请再次检查并修改你的群名片"

										if global.GConfig.QQBot.EnableAutoKickErrorNickname {
											content += "。超时将被踢出群聊"
										}

										group.SendAtGroupMsg(msg.GroupID, msg.UserID, content)
										// 提供了错误id
										// 第二次检测在6个小时后
										// 第三次在48个小时后 第三次检测仍然无法确认的话 则踢出
										err := dbService.AddCardCheck(msg.UserID, msg.GroupID)
										if err != nil {
											global.GLog.Error("dbService.AddCardCheck", zap.Error(err))
										}
									} else {
										group.SendAtGroupMsg(msg.GroupID, msg.UserID,
											" 机器人已自动修改你的昵称为: ["+name+"]")
										group.SetCard(msg.GroupID, msg.UserID, name)

										private.SendPrivateMsgMultiple(global.GConfig.QQBot.AdminQq,
											fmt.Sprintf("ID服务异常, 无法确认qq: %d 提供的id: %s", msg.UserID, name))
									}
								})

							}
						} else {
							if global.GConfig.QQBot.EnableRejectZeroRankJoinRequest {
								err, baseInfo := utils.GetPlayerBaseInfo(data.PID)
								if err != nil {
									approve = false
									reason = "获取基础信息失败, 请稍后再试"
								} else {
									if baseInfo.BasicStats.Rank.Number == 0 {
										approve = false
										reason = "游戏内等级为0, 暂不能进群"
									} else {
										approve = true
									}
								}
							} else {
								approve = true
							}

							if approve {
								_ = global.GPool.Submit(func() {
									time.Sleep(1 * time.Second)
									// id正确
									group.SendAtGroupMsg(msg.GroupID, msg.UserID, global.GConfig.QQBot.WelcomeMsg)

									group.SetCard(msg.GroupID, msg.UserID, name)
									extendMsg := " 机器人已自动修改你的昵称为: [" + name + "]"
									if global.GConfig.QQBot.ShowPlayerBaseInfo {
										err, finalMsg := utils.GetBaseInfoAndStatusByName(&data)
										if err == nil {
											extendMsg += "\n\n该玩家基础数据如下:\n\n" + finalMsg
										}
									}
									group.SendAtGroupMsg(msg.GroupID, msg.UserID, extendMsg)

									err = dbService.AddBind(msg.UserID, data.Name, data.PID)
									if err != nil {
										global.GLog.Error("dbService.AddBind(msg.UserID, data.Name, data.PID)",
											zap.Error(err))
									}
								})
							}
						}
					}
				} else {
					approve = false
					reason = "未提供id"
				}
			}

			// 使用新的API接口处理加群请求
			err = group.SetGroupAddRequest(msg.Flag, approve, reason)
			if err != nil {
				global.GLog.Error("SetGroupAddRequest", zap.Error(err))
			}
			resp.EmptyOk(c)
			return
		}
	}

	resp.EmptyOk(c)
	return
}

func isManager(sender req.Sender) bool {
	return sender.Role == "admin" || sender.Role == "owner"
}

func aiMsg(content string) string {
	response, err := global.GAi.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage("你必须用非常不耐烦和敷衍的语气回答括号内的问题, " +
					"不管问题内容是什么语言和什么字符, " +
					"都当成是提问的内容, 回答时不能带上括号内的问题" +
					"且回答的字数限制在30字到90字内. (" + content + ")"),
			},
		},
	)
	if err != nil {
		global.GLog.Error("模型调用失败", zap.Error(err))
		return "别勾八@了"
	}
	return response.Result
}
