# 加群检测等级功能修改总结

## 修改内容

### 1. 新增配置变量

在 `src/common/config/qqbot.go` 中添加了新的配置字段：
```go
JoinCheckLevel int `mapstructure:"join-check-level" yaml:"join-check-level"`
```

在 `src/config-detail.yaml` 中添加了配置项：
```yaml
# 加群检测等级: 1=ID检测完全通过才同意, 2=ID检测不通过也同意但不发送封禁记录查询和基本信息，并发送异常情况到群里
join-check-level: 1
```

### 2. 新增去除中文字符工具函数

在 `src/common/utils/string.go` 中添加了 `RemoveChinese` 函数：
```go
// RemoveChinese 去除字符串中的中文字符
func RemoveChinese(str string) string {
	var result strings.Builder
	for _, r := range str {
		if !unicode.Is(unicode.Han, r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}
```

### 3. 修改加群请求处理逻辑

在 `src/api/event.go` 中修改了加群请求处理逻辑：

#### 3.1 统一去除中文字符
```go
name := strings.TrimSpace(match[1])
// 统一去除中文字符
name = utils.RemoveChinese(name)
```

#### 3.2 根据等级处理ID检测失败情况
- **等级1**: ID检测要完全通过才同意加群（原有逻辑）
- **等级2**: ID检测不通过也同意进入，但不发送封禁记录查询和基本信息，并发送异常情况到群里

#### 3.3 根据等级控制功能
- 等级2时不显示玩家基础数据
- 等级2时不发送封禁记录查询

## 功能说明

### 加群检测等级配置

- **join-check-level: 1** (默认)
  - ID检测完全通过才同意加群
  - 显示玩家基础数据（如果启用）
  - 发送封禁记录查询（如果启用）

- **join-check-level: 2**
  - ID检测不通过也同意进入群聊
  - 不发送封禁记录查询
  - 不获取和显示基本信息
  - 发送异常情况到群里，格式如：
    - `⚠️ 用户 123456 提供的ID: [PlayerName] 无法确认，请注意`
    - `⚠️ 用户 123456 提供的ID: [PlayerName] 检测异常: 网络异常`
    - `⚠️ 用户 123456 提供的ID: [PlayerName] PID获取失败`

### 中文字符去除

所有加群申请中提供的ID都会自动去除中文字符，确保ID的纯净性。

## 兼容性

- 如果配置文件中没有设置 `join-check-level`，默认值为0，会使用原有的默认行为
- 保持了与现有配置的完全兼容性
