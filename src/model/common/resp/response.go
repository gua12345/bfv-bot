package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmptyResponse struct {
}

type Data struct {
	File    string `json:"file"`
	Summary string `json:"summary"`
}

type ImageMessage struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type ImageReply struct {
	Reply []ImageMessage `json:"reply"`
}

type ReplyResponse struct {
	Reply string `json:"reply"`
}

// 新增：回复消息段数据结构
type ReplyData struct {
	Id interface{} `json:"id"` // 支持string和int64
}

type ReplyMessage struct {
	Type string    `json:"type"`
	Data ReplyData `json:"data"`
}

type TextData struct {
	Text string `json:"text"`
}

type TextMessage struct {
	Type string   `json:"type"`
	Data TextData `json:"data"`
}

// 新增：支持回复的消息响应结构
type ReplyWithMessageResponse struct {
	Reply []interface{} `json:"reply"`
}

func EmptyOk(c *gin.Context) {
	c.JSON(http.StatusOK, EmptyResponse{})
}

func ReplyOk(c *gin.Context, reply string) {
	c.JSON(http.StatusOK, ReplyResponse{reply})
}

func ReplyWithData(c *gin.Context, m map[string]interface{}) {
	c.JSON(http.StatusOK, m)
}

func ImageOk(c *gin.Context, path string, nickname string) {
	data := make([]ImageMessage, 1)
	data[0] = ImageMessage{Type: "image", Data: Data{
		File:    path,
		Summary: nickname,
	}}
	c.JSON(http.StatusOK, ImageReply{Reply: data})
}

// 新增：支持回复消息的文本响应
func ReplyWithReply(c *gin.Context, messageId int64, message string) {
	reply := []interface{}{
		ReplyMessage{
			Type: "reply",
			Data: ReplyData{Id: messageId},
		},
		TextMessage{
			Type: "text",
			Data: TextData{Text: message},
		},
	}
	c.JSON(http.StatusOK, ReplyWithMessageResponse{Reply: reply})
}

// 新增：支持回复消息的图片响应
func ImageWithReply(c *gin.Context, messageId int64, file string, summary string) {
	reply := []interface{}{
		ReplyMessage{
			Type: "reply",
			Data: ReplyData{Id: messageId},
		},
		ImageMessage{
			Type: "image",
			Data: Data{File: file, Summary: summary},
		},
	}
	c.JSON(http.StatusOK, ReplyWithMessageResponse{Reply: reply})
}
