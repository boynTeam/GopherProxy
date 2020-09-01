package pkg

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
	DefaultCode    = 200
	DefaultMessage = "succ"

	RateLimitErrorCode        = 1000
	DefaultRateLimitErrorCode = "rate limit not allow"

	CircuitBreakerErrorCode = 1001

	DownstreamErrorCode = 1002
)

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

func (b *MessageBuilder) Message(message string) *MessageBuilder {
	b.m.Message = message
	return b
}

func (b *MessageBuilder) Data(data interface{}) *MessageBuilder {
	b.m.Data = data
	return b
}

func (b *MessageBuilder) Build() Message {
	return b.m
}
