package dto

type GetGroupMemberInfoResp struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Wording string                 `json:"wording"`
	Retcode int                    `json:"retcode"`
	Data    GetGroupMemberInfoData `json:"data,omitempty"`
}

type GetGroupMemberInfoData struct {
	GroupId int64  `json:"group_id"`
	UserId  int64  `json:"user_id"`
	Card    string `json:"card"`
}

type GetGroupMemberListResp struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Wording string                   `json:"wording"`
	Retcode int                      `json:"retcode"`
	Data    []GetGroupMemberListData `json:"data,omitempty"`
}

// 新增：获取群列表相关结构
type GetGroupListResp struct {
	Status  string             `json:"status"`
	Retcode int                `json:"retcode"`
	Data    []GetGroupListData `json:"data"`
}

type GetGroupListData struct {
	GroupId   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}

// 新增：获取消息相关结构
type GetMsgResp struct {
	Status  string     `json:"status"`
	Retcode int        `json:"retcode"`
	Data    GetMsgData `json:"data"`
}

type GetMsgData struct {
	Time        int64       `json:"time"`
	MessageType string      `json:"message_type"`
	MessageId   int64       `json:"message_id"`
	RealId      int64       `json:"real_id"`
	Sender      MsgSender   `json:"sender"`
	Message     interface{} `json:"message"`
}

type MsgSender struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
	Area     string `json:"area"`
	Level    string `json:"level"`
	Role     string `json:"role"`
	Title    string `json:"title"`
}

// 新增：获取群历史聊天记录相关结构
type GetGroupMsgHistoryResp struct {
	Status  string                  `json:"status"`
	Retcode int                     `json:"retcode"`
	Data    GetGroupMsgHistoryData  `json:"data"`
}

type GetGroupMsgHistoryData struct {
	Messages []GroupHistoryMessage `json:"messages"`
}

type GroupHistoryMessage struct {
	MessageType string      `json:"message_type"`
	SubType     string      `json:"sub_type"`
	MessageId   int64       `json:"message_id"`
	GroupId     int64       `json:"group_id"`
	UserId      int64       `json:"user_id"`
	Anonymous   interface{} `json:"anonymous"`
	Message     interface{} `json:"message"`
	RawMessage  string      `json:"raw_message"`
	Font        string      `json:"font"`
	Sender      MsgSender   `json:"sender"`
}

// 新增：获取群信息相关结构
type GetGroupInfoResp struct {
	Status  string          `json:"status"`
	Retcode int             `json:"retcode"`
	Data    GetGroupInfoData `json:"data"`
}

type GetGroupInfoData struct {
	GroupId        int64  `json:"group_id"`
	GroupName      string `json:"group_name"`
	MemberCount    int    `json:"member_count"`
	MaxMemberCount int    `json:"max_member_count"`
}

// 新增：处理加群请求相关结构
type SetGroupAddRequestResp struct {
	Status  string      `json:"status"`
	Retcode int         `json:"retcode"`
	Data    interface{} `json:"data"`
}

type GetGroupMemberListData struct {
	GroupId int64  `json:"group_id"`
	UserId  int64  `json:"user_id"`
	Card    string `json:"card"`
}
