package pkg

import (
	"fmt"
)

// shallow wrap for gin response
// Author:Boyn
// Date:2020/9/1

type Message struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type MessageBuilder struct {
	m Message
}

const (
	DefaultCode         = 200
	DefaultMessage      = "succ"
	DefaultErrorMessage = "error"

	RateLimitErrorCode = 1000

	CircuitBreakerErrorCode = 1001

	DownstreamErrorCode = 1002

	ParamErrorCode = 2000
	DbErrorCode    = 2001
)

func ErrorMessage(code int, message ...string) Message {
	var m string
	if len(message) == 0 {
		m = DefaultMessage
	} else {
		m = message[0]
	}
	return (&MessageBuilder{m: Message{
		Code:    code,
		Message: m,
	}}).Build()
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{m: Message{
		Code:    DefaultCode,
		Message: DefaultMessage,
	}}
}

func (b *MessageBuilder) Code(c int) *MessageBuilder {
	b.m.Code = c
	return b
}

func (b *MessageBuilder) Message(message interface{}) *MessageBuilder {
	b.m.Message = fmt.Sprintf("%+v", message)
	return b
}

func (b *MessageBuilder) Data(data interface{}) *MessageBuilder {
	b.m.Data = data
	return b
}

func (b *MessageBuilder) Build() Message {
	return b.m
}
