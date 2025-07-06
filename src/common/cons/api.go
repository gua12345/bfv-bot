package cons

import "bfv-bot/common/global"

// API URL 变量，在程序启动时从配置文件初始化
var (
	// BfBanStatus bfban状态
	BfBanStatus string

	// BfvRobotStatus 机器人社区状态
	BfvRobotStatus string

	// PlayerData 玩家生涯数据
	PlayerData string

	// ActiveTag 代表战排
	ActiveTag string

	// CheckPlayer 查询玩家pid/userid
	CheckPlayer string

	// Captcha 验证码
	Captcha string

	// BanLog 屏蔽记录
	BanLog string

	// ServerListBfvRobot 服务器列表
	ServerListBfvRobot string

	// ServerListGameTools 服务器列表
	ServerListGameTools string

	// Ban 屏蔽
	Ban string

	// RemoveBan 解除屏蔽
	RemoveBan string

	// ServerPlayerGameTools 服内玩家
	ServerPlayerGameTools string

	// ServerPlayerBfvRobot 服内玩家
	ServerPlayerBfvRobot string

	// GameTof 游戏任务数据
	GameTof string

	// JoinPlatoons 玩家加入的战排
	JoinPlatoons string

	// PlayerBaseInfo 玩家基础数据
	PlayerBaseInfo string

	// BfBanBatchStatus 批量获取
	BfBanBatchStatus string

	// BfvRobotBatchStatus 批量获取
	BfvRobotBatchStatus string

	// GameToolsBatchStatus gt的批量接口
	GameToolsBatchStatus string
)

// InitApiUrls 初始化API URL，从配置文件中读取baseurl并设置完整的API地址
func InitApiUrls() {
	config := global.GConfig.Api

	// 设置默认值
	bfbanBase := "https://api.bfban.com"
	bfvrobotBase := "https://api.bfvrobot.net"
	gametoolsBase := "https://api.gametools.network"

	// 如果配置文件中有自定义的baseurl，则使用配置的值
	if config.BfBanBaseUrl != "" {
		bfbanBase = config.BfBanBaseUrl
	}
	if config.BfvRobotBaseUrl != "" {
		bfvrobotBase = config.BfvRobotBaseUrl
	}
	if config.GameToolsBaseUrl != "" {
		gametoolsBase = config.GameToolsBaseUrl
	}

	// 初始化所有API URL
	BfBanStatus = bfbanBase + "/api/player"
	BfBanBatchStatus = bfbanBase + "/api/player/batch"

	BfvRobotStatus = bfvrobotBase + "/api/player/getCommunityStatus"
	PlayerData = bfvrobotBase + "/api/worker/player/getAllStats"
	ActiveTag = bfvrobotBase + "/api/worker/platoon/getActiveTags"
	CheckPlayer = bfvrobotBase + "/api/bfv/player"
	Captcha = bfvrobotBase + "/api/captcha"
	BanLog = bfvrobotBase + "/api/player/getBannedLogsByPersonaId"
	ServerListBfvRobot = bfvrobotBase + "/api/bfv/servers"
	Ban = bfvrobotBase + "/api/server/ban"
	RemoveBan = bfvrobotBase + "/api/server/removeban"
	ServerPlayerBfvRobot = bfvrobotBase + "/api/bfv/players"
	GameTof = bfvrobotBase + "/api/worker/getTOF"
	JoinPlatoons = bfvrobotBase + "/api/worker/platoon/getPlayerJoins"
	PlayerBaseInfo = bfvrobotBase + "/api/worker/player/getStats"
	BfvRobotBatchStatus = bfvrobotBase + "/api/worker/player/getBatchAllStats"

	ServerListGameTools = gametoolsBase + "/bfv/servers/"
	ServerPlayerGameTools = gametoolsBase + "/bfv/players/"
	GameToolsBatchStatus = gametoolsBase + "/bfv/multiple/"
}
