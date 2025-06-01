package group

import (
	"bfv-bot/common/des"
	"bfv-bot/common/global"
	"bfv-bot/common/http"
	"bfv-bot/model/dto"
	"errors"
	"go.uber.org/zap"
)

func SetGroupKick(groupId int64, userId int64) {
	data := map[string]interface{}{
		"group_id":           groupId,
		"user_id":            userId,
		"reject_add_request": false,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/set_group_kick", data)
}

func GetGroupMemberInfo(groupId int64, userId int64) (error, dto.GetGroupMemberInfoData) {
	data := map[string]interface{}{
		"group_id": groupId,
		"user_id":  userId,
		"no_cache": true,
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/get_group_member_info", data)
	if err != nil {
		global.GLog.Error("get_group_member_info", zap.Error(err))
		return err, dto.GetGroupMemberInfoData{}
	}
	var result dto.GetGroupMemberInfoResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, dto.GetGroupMemberInfoData{}
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("get_group_member_info", zap.String("message", result.Message),
			zap.String("wording", result.Wording))
		return errors.New("bot接口异常"), dto.GetGroupMemberInfoData{}
	}
	return nil, result.Data
}

func SetGroupWholeBan(groupId int64, enable bool) {
	data := map[string]interface{}{
		"group_id": groupId,
		"enable":   enable,
	}
	_, _ = http.Post(global.GConfig.QQBot.Address+"/set_group_whole_ban", data)
}

func GetGroupMemberList(groupId int64) (error, []dto.GetGroupMemberListData) {
	data := map[string]interface{}{
		"group_id": groupId,
		"no_cache": true,
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/get_group_member_list", data)
	if err != nil {
		global.GLog.Error("get_group_member_list", zap.Error(err))
		return err, []dto.GetGroupMemberListData{}
	}
	var result dto.GetGroupMemberListResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, []dto.GetGroupMemberListData{}
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("get_group_member_list", zap.String("message", result.Message),
			zap.String("wording", result.Wording))

		return errors.New("bot接口异常"), []dto.GetGroupMemberListData{}
	}
	return nil, result.Data
}

// GetActiveGroupMemberCardMap 获取所有启用机器人服务群的群成员名片
func GetActiveGroupMemberCardMap() map[string]bool {
	memberMap := make(map[string]bool)
	for _, item := range global.GConfig.QQBot.ActiveGroup {
		err, memberList := GetGroupMemberList(item)
		if err != nil {
			continue
		}
		for _, member := range memberList {
			if member.Card == "" {
				continue
			}
			memberMap[member.Card] = true
		}
	}
	return memberMap
}

// 新增：获取群列表
func GetGroupList(noCache bool) (error, []dto.GetGroupListData) {
	data := map[string]interface{}{
		"no_cache": noCache,
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/get_group_list", data)
	if err != nil {
		global.GLog.Error("get_group_list", zap.Error(err))
		return err, []dto.GetGroupListData{}
	}
	var result dto.GetGroupListResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, []dto.GetGroupListData{}
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("get_group_list", zap.String("status", result.Status), zap.Int("retcode", result.Retcode))
		return errors.New("bot接口异常"), []dto.GetGroupListData{}
	}
	return nil, result.Data
}

// 新增：获取消息
func GetMsg(messageId int64) (error, dto.GetMsgData) {
	data := map[string]interface{}{
		"message_id": messageId,
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/get_msg", data)
	if err != nil {
		global.GLog.Error("get_msg", zap.Error(err))
		return err, dto.GetMsgData{}
	}
	var result dto.GetMsgResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, dto.GetMsgData{}
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("get_msg", zap.String("status", result.Status), zap.Int("retcode", result.Retcode))
		return errors.New("bot接口异常"), dto.GetMsgData{}
	}
	return nil, result.Data
}

// 新增：获取群历史聊天记录
func GetGroupMsgHistory(groupId int64, messageId string, count int) (error, dto.GetGroupMsgHistoryData) {
	data := map[string]interface{}{
		"group_id":   groupId,
		"message_id": messageId,
		"count":      count,
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/get_group_msg_history", data)
	if err != nil {
		global.GLog.Error("get_group_msg_history", zap.Error(err))
		return err, dto.GetGroupMsgHistoryData{}
	}
	var result dto.GetGroupMsgHistoryResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, dto.GetGroupMsgHistoryData{}
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("get_group_msg_history", zap.String("status", result.Status), zap.Int("retcode", result.Retcode))
		return errors.New("bot接口异常"), dto.GetGroupMsgHistoryData{}
	}
	return nil, result.Data
}

// 新增：获取群信息
func GetGroupInfo(groupId int64, noCache bool) (error, dto.GetGroupInfoData) {
	data := map[string]interface{}{
		"group_id": groupId,
		"no_cache": noCache,
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/get_group_info", data)
	if err != nil {
		global.GLog.Error("get_group_info", zap.Error(err))
		return err, dto.GetGroupInfoData{}
	}
	var result dto.GetGroupInfoResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err, dto.GetGroupInfoData{}
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("get_group_info", zap.String("status", result.Status), zap.Int("retcode", result.Retcode))
		return errors.New("bot接口异常"), dto.GetGroupInfoData{}
	}
	return nil, result.Data
}

// 新增：处理加群请求/邀请
func SetGroupAddRequest(flag string, approve bool, reason string) error {
	data := map[string]interface{}{
		"flag":    flag,
		"approve": approve,
	}
	if reason != "" {
		data["reason"] = reason
	}
	post, err := http.Post(global.GConfig.QQBot.Address+"/set_group_add_request", data)
	if err != nil {
		global.GLog.Error("set_group_add_request", zap.Error(err))
		return err
	}
	var result dto.SetGroupAddRequestResp
	err = des.StringToStruct(post, &result)
	if err != nil {
		global.GLog.Error("StringToStruct", zap.Error(err))
		return err
	}
	if result.Status != "ok" || result.Retcode != 0 {
		global.GLog.Error("set_group_add_request", zap.String("status", result.Status), zap.Int("retcode", result.Retcode))
		return errors.New("bot接口异常")
	}
	return nil
}
