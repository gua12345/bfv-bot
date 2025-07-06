package config

// Api API相关配置
type Api struct {
	// BfBan API 基础URL
	BfBanBaseUrl string `mapstructure:"bfban-base-url" yaml:"bfban-base-url"`
	// BfvRobot API 基础URL
	BfvRobotBaseUrl string `mapstructure:"bfvrobot-base-url" yaml:"bfvrobot-base-url"`
	// GameTools API 基础URL
	GameToolsBaseUrl string `mapstructure:"gametools-base-url" yaml:"gametools-base-url"`
}
