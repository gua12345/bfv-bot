package flow

import (
	"bfv-bot/bot/group"
	"bfv-bot/bot/private"
	"bfv-bot/common/global"
	"time"
)

func CleanExpiredPrivateFlow() {

	deleteList := make([]int64, 0)

	now := time.Now().UnixMilli()
	for key, value := range PrivateFlowable {
		if now-(60*1000) > value.ActiveTime {
			private.SendPrivateReplyMsg(key, value.MsgId, "当前对话已超时, 请重新发起")
			deleteList = append(deleteList, key)
		}
	}

	for _, item := range deleteList {
		delete(PrivateFlowable, item)
	}

}

func CleanExpiredGroupFlow() {

	deleteList := make([]string, 0)

	now := time.Now().UnixMilli()
	for key, value := range GroupFlowable {
		if now-(60*1000) > value.ActiveTime {
			group.SendGroupReplyMsg(value.GroupId, value.MsgId, "当前对话已超时, 请重新发起")
			deleteList = append(deleteList, key)
		}
	}

	for _, item := range deleteList {
		DeleteGroupStepByKey(item)
	}

}

func CurfewLog(isStart bool) {
	if isStart {
		global.GLog.Info("宵禁时间已到，机器人进入宵禁模式。")
	} else {
		global.GLog.Info("宵禁时间结束，机器人解除宵禁模式。")
	}
}
